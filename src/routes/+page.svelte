<script>
  import bg from "$lib/images/bg.jpg";
  // import goto from "$app/navigation";

  const handleClick = () => {
      // goto("/chat");
  };
</script>

<div id="body" style="background-image: url({bg}); height: 1000px;" class="relative overflow-hidden p-12 text-center">
  <div id="darken"class="absolute bottom-0 left-0 right-0 top-0 h-full w-full overflow-hidden bg-fixed"
  style="background-color: rgba(0, 0, 0, 0.2)">
    <h1 class="text-7xl text-white" style="margin-top: 10%;">
        <b>Welcome to MedTalk.</b>
    </h1>
    <h4 class="text-white text-xl">
      The definite solution for understanding your medication.
    </h4>
    <button on:click={handleClick} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full" style="margin-top: 2%;">
      {"Start your MedTalk journey →"}
    </button>
  </div>
</div>

<style>

</style>