<template>
    <button v-if="!link" :class="mode">
        <slot></slot>
    </button>
    <router-link v-else :to="to" :class="mode">
        <slot></slot>
    </router-link>
</template>

<script>
export default{
    props: {
        mode: {
            type: String,
            required: false,
            default: null
        },
        link: {
            type: Boolean,
            required: false,
            default: false
        },
        to: {
            type: String,
            required: false,
            default: '/'
        }
    }
}


</script>

<style scoped>

a, button {
    text-decoration: none;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: large;
    color: var(--primary-color);
    
}

a:hover,
button:hover,
button:active{
    color: var(--text-white);
    background-color: var(--primary-color);
    border-radius: 2rem;

}

.border {
    border-radius: 2rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.26);
    padding: 0;
    margin: 2rem auto;
    max-width: 15rem;
}

</style>
