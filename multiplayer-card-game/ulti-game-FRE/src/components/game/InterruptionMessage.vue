<template>
    <div v-if="show" class="interruption-overlay">
        <div class="interruption-message">
            <h2>{{ message }}</h2>
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
        message: {
            type: String,
            default: 'A player has left the game. Returning to lobby...'
        }
    }
}
</script>

<style scoped>
.interruption-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
}

.interruption-message {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    color: white;
    padding: 3rem 4rem;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
    animation: slideIn 0.3s ease-out;
}

.interruption-message h2 {
    margin: 0;
    font-size: 2rem;
    text-align: center;
}

@keyframes slideIn {
    from {
        transform: translateY(-20px);
        opacity: 0;
    }
    to {
        transform: translateY(0);
        opacity: 1;
    }
}
</style>
