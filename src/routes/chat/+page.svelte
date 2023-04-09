<script>
    import Message from "$lib/components/Message.svelte";
    import MessageClient from "$lib/data/client.js";
    import Examples from "../../lib/components/Examples.svelte";

    let userInput = "";
    let stream = [];

    let messages = {};

    const handleSubmit = () => {
        let history = [];
        for (let i = 0; i < stream.length / 2; i++) {
            let obj = {
                "human": stream[i * 2].trim(),
                "assistant": stream[i * 2 + 1].trim()
            }
            history.push(obj);
        }
        
        // console.log(userInput);
        const client = new MessageClient('https://conversation-api-test.yarn.network/conversation', userInput.trim(), history);
        stream = [...stream, userInput, ""];
        client.connect();
        client.on('newtext', (text) => {
            console.log('Received new text:', text);
            stream[stream.length - 1] += text;
        });
        client.on('finaltext', (text, links) => {
            console.log('Received final:', text, links)
            stream = [...stream.slice(0, stream.length - 1), text]; // append links[]
            messages[stream.length-1] = {
                text, links
            }
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
                top: document.body.scrollHeight,
                behavior: "auto",
            });
        }
    };
</script>

<div class="flex flex-col px-24 py-12">
    {#if stream.length === 0}     
        <h1 class="text-5xl text-center">MedTalk Chat</h1>
    {/if}
    <div class="w-full flex flex-col gap-5 mb-12">
        {#if stream.length === 0}
            <Examples onPick={(prompt) => {
                userInput = prompt
                handleSubmit()
            }}/>
        {/if}
        {#each stream as message, i}
            <Message text={message} message={messages[i] ?? {}} count={i}/>
        {/each}
    </div>
</div>

<form
    on:submit|preventDefault={handleSubmit}
    class="fixed bottom-8 flex flex-row gap-2 justify-center"
>
    <button
        class="border border-teal-900 rounded px-5 py-2 bg-white text-teal-800"
        on:click={() => {
            location.reload()
        }}
    >
    ♻️
    </button>
    <textarea
        rows="1"
        placeholder="Send a message..."
        class="w-96 resize-none px-4 py-2 border border-teal-900 rounded text-teal-800"
        bind:value={userInput}
        on:keypress={keyHandler}
    />
    <button
        id="submitBtn"
        type="submit"
        class="border border-teal-900 rounded px-5 py-2 bg-white text-teal-800"
    >
        Enter
    </button>
</form>
