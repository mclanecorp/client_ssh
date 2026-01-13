# MySSH GUI

Interface graphique pour le client SSH/SCP MySSH, construite avec Tauri et Svelte.

## Prérequis

### Dépendances système (Linux/Debian)

```bash
sudo apt-get update
sudo apt-get install -y \
  libwebkit2gtk-4.1-dev \
  libgtk-3-dev \
  libayatana-appindicator3-dev \
  librsvg2-dev \
  patchelf \
  libssl-dev \
  build-essential \
  curl \
  wget \
  file
```

### Outils de développement

- Node.js >= 18
- Rust et Cargo
- Le binaire `myssh` compilé (dans le répertoire parent)

## Installation

```bash
# Installer les dépendances npm
npm install
```

## Développement

Pour lancer l'application en mode développement :

```bash
npm run tauri dev
```

L'application se lancera et rechargera automatiquement lors des modifications du code frontend.

**Note:** L'application utilise le binaire `myssh` situé dans `../myssh` (répertoire parent). Assurez-vous qu'il existe et est exécutable.

## Build

Pour créer un build de production :

```bash
npm run tauri build
```

Le binaire sera créé dans `src-tauri/target/release/`.

## Architecture

- **Frontend**: Svelte 5 avec SvelteKit
- **Backend**: Tauri (Rust)
- **CLI Backend**: myssh (Go) - appel via `std::process::Command`

### Structure

```
myssh-gui/
├── src/                    # Code frontend Svelte
│   └── routes/
│       └── +page.svelte   # Page principale avec navigation par tabs
├── src-tauri/             # Code Rust Tauri
│   ├── src/
│   │   ├── main.rs       # Point d'entrée
│   │   └── lib.rs        # Commandes Tauri (appel vers myssh CLI)
│   └── tauri.conf.json   # Configuration Tauri
└── static/                # Assets statiques
```

### Commandes Tauri disponibles

L'application expose ces commandes Tauri qui appellent le CLI myssh :

- `ssh_connect` - Connexion SSH interactive
- `scp_upload` - Upload de fichier via SCP
- `scp_download` - Download de fichier via SCP
- `profile_list` - Liste les profils enregistrés
- `profile_add` - Ajoute un nouveau profil
- `profile_delete` - Supprime un profil

## Configuration

L'application cherche le binaire `myssh` à l'emplacement suivant :

1. Variable d'environnement `MYSSH_PATH` si définie
2. Sinon, `../myssh` (répertoire parent)

Pour utiliser un chemin personnalisé :

```bash
export MYSSH_PATH=/path/to/myssh
npm run tauri dev
```

## Fonctionnalités

### 1. Connexion SSH

- Connexion via profil enregistré ou paramètres directs (host, user, port)
- Support de l'authentification par mot de passe et/ou clé privée
- Session SSH interactive

### 2. Transfert SCP

- Upload et download de fichiers
- Support des profils ou paramètres directs
- Feedback en temps réel

### 3. Gestion des profils

- Ajouter, lister et supprimer des profils SSH
- Stockage dans SQLite (~/.config/myssh/myssh.db)
- Réutilisation facile des profils dans SSH et SCP

## Développement

### Variables d'environnement

- `MYSSH_PATH` - Chemin vers le binaire myssh (défaut: ../myssh)

### Mode sombre

L'application supporte automatiquement le mode sombre en fonction des préférences système.

## Troubleshooting

### "Failed to execute myssh"

Vérifiez que :
1. Le binaire `myssh` existe dans `../myssh`
2. Il est exécutable (`chmod +x ../myssh`)
3. Ou définissez `MYSSH_PATH` vers le bon chemin

### Erreur de dépendances webkit2gtk

Installez les dépendances système listées dans les prérequis ci-dessus.
