<script>
  import logo from './assets/jitterbit.png'  

  let EMPTY_ENVS = ['None']

  let project = "";
  let output = "";
  let environments = EMPTY_ENVS;
  let environment = "";
  let processing = false;

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
    processing = true;
    let result = await window.go.main.App.Extract(project, environment, output);
    processing = false;
    if (result === true) {
      project = "";
      environments = EMPTY_ENVS;
      environment = "";
      output = "";
    }
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
      <div class="flex flex-row my-6 justify-center items-center">
        {#if project === "" || environment === "" || environment === "None" || output === ""}
        <button on:click={extract} class="text-white text-2xl rounded-full text-bold bg-gray-500 px-6 pb-3 pt-2 my-3" disabled>Extract</button>  
        {:else}
        <button on:click={extract} class="text-white text-2xl rounded-full text-bold bg-[#e91889] hover:bg-[#c11472] transition duration-150 px-6 pb-3 pt-2 my-3">Extract</button>
        {/if}
        {#if processing}
        <div role="status" class="ml-4">
          <svg aria-hidden="true" class="w-8 h-8 mr-2 text-gray-200 animate-spin dark:text-gray-400 fill-[#ff902a]" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor"/>
            <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill"/>
          </svg>
          <span class="sr-only">Loading...</span>
        </div>
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
