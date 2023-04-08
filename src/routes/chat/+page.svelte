<script>

    const handleSubmit = () => {
        console.log("enter");
    }
</script>

<h1 class="text-5xl">
    Chat
</h1>


<form on:submit|preventDefault={handleSubmit}>
    <input type="text" class="">
    <button>Enter</button>
</form>

