<template>
  <div class="max-w-7xl mx-auto px-12">
    <div class="mb-10">
      <h2 class="text-3xl font-bold text-slate-900 mb-2">Settings</h2>
      <p class="text-slate-600 font-medium">Configure application preferences and options</p>
    </div>

    <div class="bg-white rounded-xl border-2 border-slate-200 shadow-sm p-8">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
          <i class="fas fa-cog text-blue-600 text-xl"></i>
        </div>
        <div>
          <h3 class="text-xl font-bold text-slate-900">Application Settings</h3>
          <p class="text-sm text-slate-600">Manage your preferences and configuration</p>
        </div>
      </div>

      <div class="space-y-5">
        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-bell text-slate-600"></i>
            <div>
              <span class="font-semibold text-slate-700 block">Notifications</span>
              <span class="text-xs text-slate-500">Manage notification preferences</span>
            </div>
          </div>
          <button class="text-sm font-medium text-blue-600 hover:text-blue-700">
            Configure
          </button>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-palette text-slate-600"></i>
            <div>
              <span class="font-semibold text-slate-700 block">Theme</span>
              <span class="text-xs text-slate-500">Customize appearance</span>
            </div>
          </div>
          <button class="text-sm font-medium text-blue-600 hover:text-blue-700">
            Customize
          </button>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-language text-slate-600"></i>
            <div>
              <span class="font-semibold text-slate-700 block">Language</span>
              <span class="text-xs text-slate-500">Select your preferred language</span>
            </div>
          </div>
          <select class="text-sm font-medium text-slate-900 bg-white border border-slate-300 rounded px-3 py-1">
            <option>English</option>
            <option>Spanish</option>
            <option>Portuguese</option>
          </select>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-clock text-slate-600"></i>
            <div>
              <span class="font-semibold text-slate-700 block">Timezone</span>
              <span class="text-xs text-slate-500">Set your timezone</span>
            </div>
          </div>
          <select class="text-sm font-medium text-slate-900 bg-white border border-slate-300 rounded px-3 py-1">
            <option>UTC</option>
            <option>America/New_York</option>
            <option>America/Los_Angeles</option>
            <option>Europe/London</option>
          </select>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-shield-alt text-slate-600"></i>
            <div>
              <span class="font-semibold text-slate-700 block">Security</span>
              <span class="text-xs text-slate-500">Manage security settings</span>
            </div>
          </div>
          <button class="text-sm font-medium text-blue-600 hover:text-blue-700">
            Manage
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Settings view component
// Future: Add settings management logic here
</script>
