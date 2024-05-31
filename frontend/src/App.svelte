<script>
  import bkg from "./assets/images/bkg.png";
  import logo from "./assets/images/logo.png";
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
      <div class="h-3"></div>
      <p class="mb-5">
        Custom moded minecraft client tailored for best experience<br /> in use
        with my special servers :)<br /><br />
      </p>
      <input
        style="margin-right: 1em; --wails-draggable:no-drag"
        type="nick"
        placeholder="nickname"
        class="input input-bordered w-64 max-w-xs"
      />
      <button style="--wails-draggable:no-drag" class="btn btn-primary"
        >Launch Game</button
      >
      <div class="h-3"></div>
    </div>
  </div>
</div>
