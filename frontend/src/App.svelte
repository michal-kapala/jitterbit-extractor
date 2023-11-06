<script>
  import logo from './assets/jitterbit.png'

  let EMPTY_ENVS = ['None']

  let project = "";
  let output = "";
  let environments = EMPTY_ENVS;
  let environment = "";

  async function selectProject() {
    project = await window.go.main.App.SelectProject();
    let result = await window.go.main.App.GetEnvs(project);
    if (result !== null && result !== undefined && result != [])
      environments = result;
    else
      environments = EMPTY_ENVS;
    environment = "";
  }

  async function selectOutput() {
    output = await window.go.main.App.SelectOutput();
  }

  async function extract() {
    
  }
</script>

<main class="flex flex-row p-4 bg-gradient-to-br from-[#0f0225]/[.95] via-[#6020d6] to-[#ffaf44] text-gray-300 justify-center items-center">
  <div class="flex flex-row w-full py-64 justify-center items-center"> 

    <img src={logo} alt="Jitterbit Logo" class="m-6"/>

    <div class="flex flex-col ml-6 my-4 w-1/2">
      <div class="my-2">
        <p class="bold py-2 text-bold text-xl">Jitterbit project</p>
        <div id="input" data-wails-no-drag class="flex flex-row items-center w-full">
          <button on:click={selectProject} class="text-white rounded-full text-bold bg-[#ff902a] hover:bg-[#f67600] transition duration-150 px-3 py-2 my-2">Select</button>
          <input bind:value={project} placeholder="None" class="text-black flex ml-4 p-2 border-2 border-black rounded border-1 bg-gray-300 truncate w-full" readonly>
        </div>
      </div>
      <div class="my-2">
        <p class="bold py-2 text-bold text-xl">Environment</p>
        <div id="input" data-wails-no-drag class="flex flex-row items-center w-full">
          <select bind:value={environment} name="envs" placeholder="None" class="text-black flex p-2 border-2 border-black rounded border-1 bg-gray-300 min-w-[25%] w-fit">
            {#each environments as env}
            <option value={env}>{env}</option>
            {/each}
          </select>
        </div>
      </div>
      <div class="my-2">
        <p class="bold py-2 text-bold text-xl">Output directory</p>
        <div id="input" data-wails-no-drag class="flex flex-row items-center w-full">
          <button on:click={selectOutput} class="text-white rounded-full text-bold bg-[#ff902a] hover:bg-[#f67600] transition duration-150 px-3 py-2 my-2">Select</button>
          <input bind:value={output} placeholder="None" class="text-black flex ml-4 p-2 border-2 border-black rounded border-1 bg-gray-300 truncate w-full" readonly>
        </div>
      </div>
      <div class="flex flex-col my-6 justify-center items-center">
        {#if project === "" || environment === "" || environment === "None" || output === ""}
        <button on:click={extract} class="text-white text-2xl rounded-full text-bold bg-gray-500 px-6 pb-3 pt-2 my-3" disabled>Extract</button>  
        {:else}
        <button on:click={extract} class="text-white text-2xl rounded-full text-bold bg-[#e91889] hover:bg-[#c11472] transition duration-150 px-6 pb-3 pt-2 my-3">Extract</button>
        {/if}
      </div>
    </div>

  </div>
</main>

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }

  main {
    padding: 1em;
    margin: 0 auto;
  }

  img {
    height: 16rem;
    width: 16rem;
  }
</style>
