<template>
    <button class="card-btn" @click="cardClicked" :disabled="disabled">
        <img v-if="imageLoaded" 
             :src="getCardImage()" 
             :alt="`Card ${cardId}`" 
             class="card-image"
             @error="handleImageError" />
        <span v-else>{{ cardId }}</span>
    </button>
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
    data() {
        return {
            imageLoaded: true
        };
    },
    methods: {
        cardClicked(){
            if(this.disabled){
                return;
            }
            this.$emit("clicked", this.cardId);
        },
        handleImageError() {
            this.imageLoaded = false;
        },
        getCardImage() {
            try {
                return new URL(`../../assets/kartya/${this.cardId}.jpg`, import.meta.url).href;
            } catch (e) {
                console.error(`Failed to load image for card ${this.cardId}:`, e);
                this.imageLoaded = false;
                return '';
            }
        }
    }
};
</script>

<style scoped>
.card-btn {
    padding: 0.5rem;
    border: 2px solid #0062ff;
    border-radius: 8px;
    background-color: #fff;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 120px;
    transition: transform 0.2s, box-shadow 0.2s;
}

.card-btn:hover:not(:disabled) {
    transform: translateY(-5px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.card-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.card-image {
    max-width: 100%;
    max-height: 120px;
    width: auto;
    height: auto;
    object-fit: contain;
    border-radius: 4px;
}

.card-btn span {
    font-size: 1.5rem;
    font-weight: bold;
    color: #333;
}
</style>
