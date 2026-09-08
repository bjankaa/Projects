<template>
    <div class="container">
        <div class="title-image-placeholder">
            <p>Title Image</p>
        </div>
        <div class="border">
            <router-link to="/auth">
                <h3>PLAY ULTI NOW</h3>
            </router-link>
        </div>
    </div>
</template>


<style scoped>
.container {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.title-image-placeholder {
    width: 400px;
    height: 250px;
    background: linear-gradient(135deg, #ea6666 0%, hsl(115, 37%, 46%) 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 2rem auto;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.title-image-placeholder p {
    color: white;
    font-size: 1.5rem;
    font-weight: 600;
    opacity: 0.8;
}

a {
    text-decoration: none;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: large;
    color: var(--primary-color);
    padding: 1rem 2rem;
    border-radius: 2rem;
    transition: all 0.3s ease;
}

a:hover {
    color: var(--text-white);
    background-color: var(--primary-color);
}

.border {
    border-radius: 2rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.26);
    padding: 0;
    margin: 2rem auto;
    max-width: 15rem;
    width: 15rem;
}

h3 {
    margin: 0.5rem 0;
}
</style>