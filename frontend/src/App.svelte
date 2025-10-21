<script>
  import bkg from "./assets/images/bkg.png";
  import logo from "./assets/images/logo.png";
  import { Martini } from "lucide-svelte";
  import { Cigarette } from "lucide-svelte";
  import { CookingPot } from "lucide-svelte";
  import * as runtime from "../wailsjs/runtime";
  import { Version } from '../wailsjs/go/main/App';
  import { LoadOptions } from '../wailsjs/go/main/App';
  import { SaveOptions } from '../wailsjs/go/main/App';

  let version = "";
  Version().then((v) => {
    version = v
    desc = desc + " (verzia " + version + ")";
  });

  let desc = "Vitajte v Mikiho Launcheri!";
  let blink = false;

  let nickname = "";
  LoadOptions().then((opts) => {
    if (opts && opts.nickname) {
      nickname = opts.nickname;
    }
  });

  async function start() {
    SaveOptions(nickname);
    blink = true;
    desc = "Priebieha pripájanie na server!";
  }
</script>

<div class="hero min-h-screen" style="--wails-draggable:drag">
  <img style="pointer-events: none; object-fit: cover;" src={bkg} />
  <div class="hero-overlay bg-opacity-40"></div>
  <div class="hero-content text-center text-neutral-content">
    <div class="max-w-xl">
      <img
        class="max-w-lg"
        style="pointer-events: none; object-fit: cover;"
        src={logo}
      />
      <div class="h-5"></div>
      {#if blink}
        <p class="mb-5 blink">{desc}</p>
      {:else}
        <p class="mb-5">{desc}</p>
      {/if}
      <div class="h-8"></div>
      <div class="flex justify-center">
        <input
          bind:value={nickname}
          style="margin-right: 1em; --wails-draggable:no-drag"
          type="nick"
          placeholder="nick"
          class="input input-bordered w-60 max-w-xs"
        />
        <button
          onclick="my_modal_2.showModal()"
          style="--wails-draggable:no-drag"
          class="btn btn-primary"
          ><CookingPot class="w-5" /><Martini class="w-5" /><Cigarette
            class="w-5"
          /></button
        >
        <dialog id="my_modal_2" class="modal">
          <div class="modal-box">
            <h3 style="pointer-events: none;" class="font-bold text-lg">
              Pre pokračovanie musia byť splnené tieto zásady:
            </h3>
            <br />
            <p style="pointer-events: none;">
              1. Priprav si cháles<br />
              2. Nalej si akýsi drink<br />
              3. Odpál fajnové
            </p>
            <br />
            <form method="dialog">
              <div class="flex gap-6 justify-center">
                <button
                  style="--wails-draggable:no-drag;"
                  on:click={() => start()}
                  class="btn btn-primary">Ano</button
                >
                <button
                  style="--wails-draggable:no-drag;"
                  on:click={() => runtime.Quit()}
                  class="btn btn-neutral">SKAP</button
                >
              </div>
            </form>
          </div>
          <form method="dialog" class="modal-backdrop">
            <button>close</button>
          </form>
        </dialog>
      </div>

      <div class="h-3"></div>
    </div>
  </div>
</div>

<style>
  .blink {
    animation: blink 1s linear infinite;
  }
  @keyframes blink {
    50% {
      opacity: 0;
    }
  }
</style>
