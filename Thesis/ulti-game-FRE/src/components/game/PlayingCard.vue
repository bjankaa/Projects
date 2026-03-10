<template>
    <button class="card-btn" @click="cardClicked" :disabled="disabled">{{ cardId }}</button>
</template>

<script>
export default {
    props: {
        cardId: {
            type: [Number, String],
            required: true,
        },
        disabled:{
            type: Boolean,
            default: true,
        }
    },
    methods: {
        cardClicked(){
            if(this.disabled){
                return;
            }
            this.$emit("clicked", this.cardId);
        }
    }
};
</script>

<style scoped>
.card-btn {
    padding: 1rem;
    font-size: 1.5rem;
    border: 2px solid #0062ff;
    border-radius: 8px;
    background-color: #fff;
    cursor: pointer;
}

.card-btn:hover {
    background-color: #e6e6e6;
}

.card-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}
</style>
