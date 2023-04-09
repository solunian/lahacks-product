<script>
    export let onPick;
    import Example from "./Example.svelte";
    import starSVG from "$lib/images/star.svg";

    const examples = [
        "What medication should I take when I have a high fever and a sore throat?",
        "How can I prevent myself from getting sick at work?",
        "What are some side effects of Tylenol while sleep deprived?"
    ]
</script>

<div class="flex flex-row justify-center mt-12">
    <div class="flex flex-col justify-center gap-4">
        <div class="flex flex-row justify-center">
            <img src={starSVG} width={24} height={24} alt=""/>
        </div>

        <h1 class="text-2xl text-center">
            Examples
        </h1>

        {#each examples as example}
            <Example prompt={example} onPick={() => {
                onPick(example)
            }}/>
        {/each}
    </div>
</div>

