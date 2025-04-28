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
                name: 'Home',
                redirect: { name: 'AI' },
            },
            {
                path: 'dicom',
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
    {
        path: '/viewer/:id',
        name: 'Viewer',
        component: () => import('@/views/DicomViewer.vue'),
    },
    {
        path: '/config',
        name: 'Config',
        component: () => import('@/views/Config.vue'),
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from) => {
    const hasUrlConfig = localStorage.getItem('base')
    if (hasUrlConfig || to.name === 'Config') {
        return
    }
    return { name: 'Config' }
})

export default router
