<template>
    <div class="spinner"></div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'Loading'
});
</script>

<style>
.spinner {
    width: 24px;
    height: 24px;
    display: grid;
    border: 3.5px solid #0000;
    border-radius: 50%;
    border-color: #dbdcef #0000;
    animation: spinner-e04l1k 1s infinite linear;
}

.spinner::before,
.spinner::after {
    content: '';
    grid-area: 1/1;
    margin: 2.2px;
    border: inherit;
    border-radius: 50%;
}

.spinner::before {
    border-color: #474bff #0000;
    animation: inherit;
    animation-duration: 0.5s;
    animation-direction: reverse;
}

.spinner::after {
    margin: 8.9px;
}

@keyframes spinner-e04l1k {
    100% {
        transform: rotate(0.5turn);
    }
}
</style>