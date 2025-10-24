import minecraft_launcher_lib
import subprocess
import os

MINECRAFT_VERSION = "1.20.1"                     # fixed Minecraft version
MINECRAFT_DIR = os.path.join(os.getcwd(), "mc")  # relative /mc directory
USERNAME = "misko_22"                           # any offline name
RAM = "4G"                                       # memory allocation
JAVA_PATH = os.path.join(
    MINECRAFT_DIR,
    "runtime",
    "java-runtime-gamma",
    "windows-x64",
    "java-runtime-gamma",
    "bin",
    "java.exe"
)

os.makedirs(MINECRAFT_DIR, exist_ok=True)

latest_fabric_loader = minecraft_launcher_lib.fabric.get_latest_loader_version()

# minecraft_launcher_lib.fabric.install_fabric(
#     minecraft_version=MINECRAFT_VERSION,
#     loader_version= latest_fabric_loader,
#     minecraft_directory=MINECRAFT_DIR,
#     java=JAVA_PATH
# )

version_id = f"fabric-loader-{latest_fabric_loader}-{MINECRAFT_VERSION}"

options = {
    "username": USERNAME,
    "uuid": "0",
    "token": "0",  # offline mode
    "jvmArguments": [f"-Xmx{RAM}", f"-Xms{RAM}"],
    "quickPlayMultiplayer": "mc.hicz.net",
}

cmd = minecraft_launcher_lib.command.get_minecraft_command(
    version_id, MINECRAFT_DIR, options
)



print(cmd)
#subprocess.run(cmd, cwd=MINECRAFT_DIR)
