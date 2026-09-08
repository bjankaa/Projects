<template>
    <button type="submit" :disabled="disabled" class="submit-button">
        <slot></slot>
    </button>
</template>

<script>
export default {
    name: 'SubmitButton',
    props: {
        disabled: {
            type: Boolean,
            default: false
        }
    }
}
</script>

<style scoped>
.submit-button {
    text-decoration: none;
    font-size: large;
    color: var(--text-white);
    border: 1px solid transparent;
    border-color: var(--primary-color);
    margin-top: 0.5rem;
    border-radius: 2rem;
    background: var(--primary-color);
    max-width: 15rem;
    display: inline-block;
    padding: 0.2rem 1rem;
    cursor: pointer;
    transition: all 0.3s ease;
}

.submit-button:hover:not(:disabled) {
    color: var(--primary-color);
    background-color: var(--text-white);
    border-color: var(--primary-color);
}

.submit-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}
</style>
