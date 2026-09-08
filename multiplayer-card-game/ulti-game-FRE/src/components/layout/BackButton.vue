<template>
    <button @click="goBack" class="back-button">
        <span class="arrow">←</span> Back
    </button>
</template>

<script>
export default {
    name: 'BackButton',
    methods: {
        goBack() {
            this.$router.back();
        }
    }
}
</script>

<style scoped>
.back-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background-color: transparent;
    border: 1px solid #ddd;
    border-radius: 6px;
    color: #666;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.3s ease;
    margin-top: 2rem;
}

.back-button:hover {
    background-color: #f5f5f5;
    border-color: #999;
    color: #333;
}

.arrow {
    font-size: 1.2rem;
    font-weight: bold;
}
</style>
