import minecraft_launcher_lib
import os
import subprocess

# === Configuration ===
MINECRAFT_DIR = os.path.join(os.getcwd(), "mc")  # relative directory
MINECRAFT_VERSION = "1.20.1"                     # fixed Minecraft version
JAVA_PATH = "java"                               # path to Java (must be in PATH)
RAM = "4G"                                       # memory allocation
USERNAME = "misko_22"                              # offline username

# === Ensure directory exists ===
os.makedirs(MINECRAFT_DIR, exist_ok=True)

# === Get the latest Fabric loader version for this Minecraft version ===
latest_fabric_loader = minecraft_launcher_lib.fabric.get_latest_loader_version(MINECRAFT_VERSION)
print(f"🧶 Latest Fabric loader for {MINECRAFT_VERSION}: {latest_fabric_loader}")

# === Install Fabric client ===
print("⬇️ Installing Fabric client...")
minecraft_launcher_lib.fabric.install_fabric( 
    MINECRAFT_VERSION,
    latest_fabric_loader,
    MINECRAFT_DIR
)

# Construct the version ID
version_id = f"fabric-loader-{latest_fabric_loader}-{MINECRAFT_VERSION}"
print(f"✅ Installed Fabric client: {version_id}")

# === Launch options ===
options = {
    "username": USERNAME,
    "uuid": "0",
    "token": "0",  # offline mode
    "jvmArguments": [f"-Xmx{RAM}", f"-Xms{RAM}"],
}

# === Build launch command ===
print("⚙️ Building launch command...")
cmd = minecraft_launcher_lib.command.get_minecraft_command(
    version_id, MINECRAFT_DIR, options
)

# === Launch Minecraft ===
print("🎮 Launching Fabric Minecraft client...")
subprocess.run([JAVA_PATH] + cmd)
