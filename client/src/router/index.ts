import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        redirect: '/home',
    },
    {
        path: '/home',
        component: () => import('@/layouts/default.vue'),
        children: [
            {
                path: '',
                redirect: 'dicom',
            },
            {
                path: '/dicom',
                name: 'Dicom',
                component: () => import('@/views/Dicom.vue'),
            },
            {
                path: 'ai',
                name: 'AI',
                component: () => import('@/views/Ai.vue'),
            },
        ],
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router
