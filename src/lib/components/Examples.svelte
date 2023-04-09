<script>
    import Example from "./Example.svelte";

    const examples = [
        "What medication should I take when I have a high fever and a sore throat?",
        "How can I prevent myself from getting sick at work?",
        "What are some side effects of Tylenol while sleep deprived?"
    ]
</script>

<div class="flex flex-row justify-center mt-16">
    <div class="flex flex-col justify-center gap-4">
        <h1 class="text-2xl mb-1 text-center">Examples</h1>

        {#each examples as example}
            <Example prompt={example}/>
        {/each}
    </div>
</div>

