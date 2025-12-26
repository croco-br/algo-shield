import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '@/lib/api';
import { updateVuetifyTheme } from '@/plugins/vuetify';

export interface BrandingConfig {
	id: number;
	app_name: string;
	icon_url?: string | null;
	favicon_url?: string | null;
	primary_color: string;
	secondary_color: string;
	header_color: string;
	created_at: string;
	updated_at: string;
}

export interface UpdateBrandingRequest {
	app_name: string;
	icon_url?: string | null;
	favicon_url?: string | null;
	primary_color: string;
	secondary_color: string;
	header_color: string;
}

export const useBrandingStore = defineStore('branding', () => {
	const config = ref<BrandingConfig | null>(null);
	const loading = ref(true);
	const error = ref<string | null>(null);

	async function loadBranding() {
		try {
			loading.value = true;
			error.value = null;

			// Fetch branding config from API (public endpoint, no auth required)
			const brandingData = await api.get<BrandingConfig>('/api/v1/branding');
			config.value = brandingData;

			// Apply branding immediately
			applyBranding(brandingData);
		} catch (e) {
			error.value = e instanceof Error ? e.message : 'Failed to load branding';
			console.error('Failed to load branding configuration:', e);

			// Apply default branding on error
			applyDefaultBranding();
		} finally {
			loading.value = false;
		}
	}

	async function updateBranding(data: UpdateBrandingRequest) {
		try {
			loading.value = true;
			error.value = null;

			// Update branding config via API (requires admin auth)
			const updatedConfig = await api.put<BrandingConfig>('/api/v1/branding', data);
			config.value = updatedConfig;

			// Apply new branding immediately
			applyBranding(updatedConfig);

			return updatedConfig;
		} catch (e) {
			error.value = e instanceof Error ? e.message : 'Failed to update branding';
			throw e;
		} finally {
			loading.value = false;
		}
	}

	function applyBranding(brandingConfig: BrandingConfig) {
		// Update document title
		document.title = brandingConfig.app_name;

		// Update CSS custom properties for colors
		const root = document.documentElement;
		root.style.setProperty('--color-primary', brandingConfig.primary_color);
		root.style.setProperty('--color-secondary', brandingConfig.secondary_color);
		root.style.setProperty('--color-header-background', brandingConfig.header_color);

		// Update Vuetify theme with branding colors
		updateVuetifyTheme(brandingConfig.primary_color, brandingConfig.secondary_color);

		// Update favicon if provided
		if (brandingConfig.favicon_url) {
			updateFavicon(brandingConfig.favicon_url);
		}
	}

	function applyDefaultBranding() {
		// Apply default values if API fails
		document.title = 'AlgoShield';

		const root = document.documentElement;
		root.style.setProperty('--color-primary', '#3B82F6');
		root.style.setProperty('--color-secondary', '#10B981');
		root.style.setProperty('--color-header-background', '#1e1e1e');

		// Update Vuetify theme with default colors
		updateVuetifyTheme('#3B82F6', '#10B981');
	}

	function updateFavicon(faviconUrl: string) {
		// Remove existing favicon links
		const existingLinks = document.querySelectorAll("link[rel*='icon']");
		existingLinks.forEach(link => link.remove());

		// Add new favicon link
		const link = document.createElement('link');
		link.rel = 'icon';
		link.href = faviconUrl;
		document.head.appendChild(link);
	}

	// Initialize on store creation
	// Apply default branding immediately to prevent FOUC
	applyDefaultBranding()
	loadBranding();

	return {
		config,
		loading,
		error,
		loadBranding,
		updateBranding,
	};
});
