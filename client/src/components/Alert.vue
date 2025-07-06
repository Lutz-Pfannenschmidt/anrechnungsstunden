<script setup lang="ts">
interface Props {
    type: "info" | "success" | "warning" | "error";
    soft?: boolean;
    outline?: boolean;
    dash?: boolean;
    text: string;
    reactive?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
    reactive: true
});
</script>


<template>
    <div role="alert" :class="{
        'alert': true,
        'alert-info': type == 'info',
        'alert-success': type == 'success',
        'alert-warning': type == 'warning',
        'alert-error': type == 'error',
        'alert-soft': props.soft,
        'alert-outline': props.outline,
        'alert-dash': props.dash,
        'alert-vertical': props.reactive, 'sm:alert-horizontal': props.reactive
    }">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="h-6 w-6 shrink-0 stroke-current">
            <path v-if="type == 'info'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            <path v-if="type == 'success'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            <path v-if="type == 'warning'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            <path v-if="type == 'error'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <span>{{ props.text }}</span>
        <div>
            <slot></slot>
        </div>
    </div>
</template>