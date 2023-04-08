<script>
    let userInput = "";

    const handleSubmit = () => {
        console.log(userInput);
    };

    const keyHandler = (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            document.getElementById("submitBtn").click();
            userInput = "";
        }
    };
</script>

<div class="flex flex-col px-24 py-12">
    <h1 class="text-5xl">MedTalk Chat</h1>
    <div class="w-full">
    </div>
</div>

<form
    on:submit|preventDefault={handleSubmit}
    class="fixed bottom-4 flex flex-row gap-2 justify-center"
>
    <textarea
        rows="1"
        placeholder="Send a message..."
        class="w-96 resize-none px-4 py-2 border"
        bind:value={userInput}
        on:keypress={keyHandler}
    />
    <button
        id="submitBtn"
        type="submit"
        class="border rounded px-5 py-2 bg-white"
    >
        Enter
    </button>
</form>
