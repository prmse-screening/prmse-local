<template>
    <v-card>
        <v-layout>
            <v-navigation-drawer expand-on-hover permanent rail>
                <v-list>
                    <v-list-item prepend-avatar="@/assets/logo.svg" subtitle="Screening for all." title="PRMSE" />
                </v-list>

                <v-divider />

                <v-list density="compact" nav>
                    <v-list-item prepend-icon="mdi-brain" title="AI" :to="{ name: 'AI' }" />
                    <v-list-item prepend-icon="mdi-server" title="Dicom Server" :to="{ name: 'Dicom' }" />
                </v-list>
            </v-navigation-drawer>

            <v-main>
                <router-view v-slot="{ Component }">
                    <v-fade-transition hide-on-leave>
                        <component :is="Component" />
                    </v-fade-transition>
                </router-view>
            </v-main>
        </v-layout>
    </v-card>
</template>

<script lang="ts" setup>
//
</script>
