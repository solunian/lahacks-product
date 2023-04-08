<script>

    const handleSubmit = () => {
        console.log("enter");
    }
</script>


<h1 class="text-5xl">
    Chat
</h1>

<div class="flex flex-col justify-end">
    <form on:submit|preventDefault={handleSubmit}>
        <input type="text" class="border">
        <button>Enter</button>
    </form>
</div>


