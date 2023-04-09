<script>
    import Message from "$lib/components/Message.svelte";
    import MessageClient from "$lib/data/client.js";

    let userInput = "";
    let stream = [];

    const handleSubmit = () => {
        // console.log(userInput);
        const client = new MessageClient('https://conversation-api-test.yarn.network/conversation', userInput, []);
        stream = [...stream, userInput];
        client.connect();
        client.on('newtext', (text) => console.log('Received new text:', text));
        client.on('finaltext', (text, links) => {
            console.log('Received final:', text, links)
            stream = [...stream, text]; // append links[]
        });
        client.on('error', (event) => console.log('Error:', event));
        client.on('close', () => console.log('Connection closed'));

        // Here you might want to add an error when it times out
        setTimeout(() => client.close(), 10000); // close the connection after 10 seconds
        
        userInput = "";
    };

    const keyHandler = (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            document.getElementById("submitBtn").click();
            window.scrollTo({
                top: document.body.clientHeight,
                behavior: "smooth",
            });
        }
    };
</script>

<div class="flex flex-col px-24 py-12">
    <h1 class="text-5xl text-center">MedTalk Chat</h1>
    <div class="w-full flex flex-col gap-5 mt-4 mb-12">
        {#each stream as message, i}
            <Message text={message} count={i} />
        {/each}
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
