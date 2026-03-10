<template>
    <div v-if="show" class="modal-overlay" @click.self="$emit('cancel')">
        <div class="modal-content">
            <h2>{{ title }}</h2>
            <p>{{ message }}</p>
            <div class="modal-buttons">
                <button class="btn-confirm" @click="$emit('confirm')">Exit</button>
                <button class="btn-cancel" @click="$emit('cancel')">Cancel</button>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        show: {
            type: Boolean,
            required: true
        },
        title: {
            type: String,
            default: 'Exit Game?'
        },
        message: {
            type: String,
            default: 'Are you sure you want to exit? This will end the game for all players.'
        }
    },
    emits: ['confirm', 'cancel']
}
</script>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.modal-content {
    background: white;
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    max-width: 400px;
    width: 90%;
}

.modal-content h2 {
    margin: 0 0 1rem 0;
    color: #333;
}

.modal-content p {
    margin: 0 0 1.5rem 0;
    color: #666;
    line-height: 1.5;
}

.modal-buttons {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
}

.modal-buttons button {
    padding: 0.5rem 1.5rem;
    border: none;
    border-radius: 6px;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
}

.btn-cancel {
    background-color: #6c757d;
    color: white;
}

.btn-cancel:hover {
    background-color: #5a6268;
}

.btn-confirm {
    background-color: #dc3545;
    color: white;
}

.btn-confirm:hover {
    background-color: #c82333;
}
</style>
