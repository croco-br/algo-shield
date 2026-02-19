<template>
    <v-navigation-drawer
        :model-value="!isMobile || isOpen"
        :width="isCollapsed ? 80 : 240"
        :temporary="isMobile"
        :permanent="!isMobile"
        :location="'left'"
        class="border-r"
        :style="{
            top: 'var(--header-height)',
            height: 'calc(100vh - var(--header-height))',
        }"
    >
        <!-- Collapse Toggle (Desktop) -->
        <template v-if="!isMobile" #append>
            <div class="d-flex justify-end pa-2">
                <v-btn
                    icon
                    variant="text"
                    size="small"
                    @click="toggleCollapse"
                    :aria-label="
                        isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'
                    "
                >
                    <v-icon>
                        {{
                            isCollapsed ? "fa-chevron-right" : "fa-chevron-left"
                        }}
                    </v-icon>
                </v-btn>
            </div>
        </template>

        <!-- Navigation Links -->
        <v-list nav density="comfortable">
            <template v-for="item in navItems" :key="item.path">
                <!-- Transactions with schema submenu -->
                <v-list-group
                    v-if="
                        item.path === '/transactions' &&
                        !isCollapsed &&
                        transactionSchemas.length > 0
                    "
                    :value="'transactions'"
                    color="primary"
                >
                    <template #activator="{ props }">
                        <v-list-item
                            v-bind="props"
                            :prepend-icon="getIcon(item.icon)"
                            :title="$t(item.label)"
                            :active="isTransactionsActive"
                            class="mx-2 mb-1"
                        />
                    </template>
                    <v-list-item
                        v-for="schema in transactionSchemas"
                        :key="schema.id"
                        :to="`/transactions?schemaId=${schema.id}`"
                        :active="isSchemaActive(schema.id)"
                        :title="`${schema.name} - ${schema.fieldCount} ${$t('common.columns')}`"
                        class="mx-2"
                        @click="isMobile && closeMobile()"
                    />
                </v-list-group>
                <!-- Regular menu items -->
                <v-list-item
                    v-else
                    :to="getNavTo(item)"
                    :active="isActive(item.path)"
                    :prepend-icon="getIcon(item.icon)"
                    :title="isCollapsed ? '' : $t(item.label)"
                    :value="item.path"
                    @click="isMobile && closeMobile()"
                    class="mx-2 mb-1"
                >
                    <template v-if="isCollapsed" #prepend>
                        <v-icon :icon="getIcon(item.icon)" />
                    </template>
                </v-list-item>
            </template>
        </v-list>
    </v-navigation-drawer>

    <!-- Mobile Overlay -->
    <v-overlay
        v-if="isMobile && isOpen"
        :model-value="isMobile && isOpen"
        class="align-start justify-start"
        @click="closeMobile"
    />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { api } from "@/lib/api";

const emit = defineEmits<{
    collapseChange: [isCollapsed: boolean];
}>();

interface NavItem {
    label: string;
    path: string;
    icon: string;
    adminOnly?: boolean;
}

interface Schema {
    id: string;
    name: string;
    extracted_fields: Array<{ path: string; type: string; nullable: boolean }>;
}

interface TransactionSchema {
    id: string;
    name: string;
    fieldCount: number;
}

const route = useRoute();
const authStore = useAuthStore();

const allNavItems: NavItem[] = [
    { label: "sidebar.dashboard", path: "/dashboard", icon: "fa-chart-line" },
    {
        label: "sidebar.transactions",
        path: "/transactions",
        icon: "fa-exchange-alt",
    },
    { label: "sidebar.schemas", path: "/schemas", icon: "fa-code" },
    { label: "sidebar.rules", path: "/rules", icon: "fa-tasks" },
    {
        label: "sidebar.permissions",
        path: "/permissions",
        icon: "fa-users-cog",
        adminOnly: true,
    },
    {
        label: "sidebar.branding",
        path: "/branding",
        icon: "fa-palette",
        adminOnly: true,
    },
];

const navItems = computed(() => {
    return allNavItems.filter((item) => !item.adminOnly || authStore.isAdmin);
});

const isCollapsed = ref(false);
const isMobile = ref(false);
const isOpen = ref(false);
const transactionSchemas = ref<TransactionSchema[]>([]);

const isTransactionsActive = computed(() => {
    return route.path === "/transactions";
});

const isSchemaActive = (schemaId: string) => {
    return route.path === "/transactions" && route.query.schemaId === schemaId;
};

const toggleCollapse = () => {
    isCollapsed.value = !isCollapsed.value;
    emit("collapseChange", isCollapsed.value);
};

const closeMobile = () => {
    isOpen.value = false;
};

const isActive = (path: string) => {
    return route.path === path || (path !== "/" && route.path.startsWith(path));
};

// Return Font Awesome icon name as-is
const getIcon = (icon: string): string => {
    return icon;
};

// When collapsed or no schemas loaded, /transactions falls into the v-else branch.
// Auto-select the first schema so the user doesn't land on an empty "No Schema" screen.
const getNavTo = (item: NavItem): string => {
    if (
        item.path === "/transactions" &&
        transactionSchemas.value.length > 0
    ) {
        return `/transactions?schemaId=${transactionSchemas.value[0]!.id}`;
    }
    return item.path;
};

const checkMobile = () => {
    isMobile.value = window.innerWidth < 960;
    if (!isMobile.value) {
        isOpen.value = false;
    }
};

async function loadTransactionSchemas() {
    // Guard: don't fetch schemas if user is not authenticated
    // This prevents 401 → refresh → 429 → force-logout loops when
    // the Sidebar briefly mounts before the router resolves to /login
    if (!authStore.user) {
        return;
    }

    try {
        const response = await api.get<{ schemas: Schema[] }>(
            "/api/v1/schemas",
        );
        transactionSchemas.value = (response?.schemas || []).map(
            (schema: Schema) => ({
                id: schema.id,
                name: schema.name,
                fieldCount: schema.extracted_fields?.length || 0,
            }),
        );
    } catch (e) {
        // Schema loading error is non-critical for sidebar
    }
}

// Load schemas when user becomes available (handles auth timing on page load)
watch(
    () => authStore.user,
    (newUser, oldUser) => {
        if (newUser && !oldUser) {
            loadTransactionSchemas();
        }
    },
    { immediate: true },
);

// Reload schemas when navigating away from /schemas (covers create/edit/delete)
watch(
    () => route.path,
    (newPath, oldPath) => {
        if (oldPath === "/schemas" && newPath !== "/schemas") {
            loadTransactionSchemas();
        }
    },
);

onMounted(() => {
    checkMobile();
    window.addEventListener("resize", checkMobile);
    emit("collapseChange", isCollapsed.value);
});

onUnmounted(() => {
    window.removeEventListener("resize", checkMobile);
});

defineExpose({
    toggleMobile: () => {
        isOpen.value = !isOpen.value;
    },
    isCollapsed,
});
</script>
