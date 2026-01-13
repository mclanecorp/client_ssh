<script>
  import { invoke } from "@tauri-apps/api/core";

  let currentTab = $state("connect");

  // SSH Connect state
  let sshHost = $state("");
  let sshUser = $state("");
  let sshPort = $state(22);
  let sshPassword = $state("");
  let sshKeyPath = $state("");
  let sshProfile = $state("");
  let connectStatus = $state("");

  // SCP state
  let scpHost = $state("");
  let scpUser = $state("");
  let scpPort = $state(22);
  let scpPassword = $state("");
  let scpKeyPath = $state("");
  let scpProfile = $state("");
  let scpLocalPath = $state("");
  let scpRemotePath = $state("");
  let scpStatus = $state("");

  // Profile state
  let profiles = $state([]);
  let profileName = $state("");
  let profileHost = $state("");
  let profileUser = $state("");
  let profilePort = $state(22);
  let profilePassword = $state("");
  let profileKeyPath = $state("");
  let profileStatus = $state("");

  async function connectSSH(event) {
    event.preventDefault();
    connectStatus = "Connexion en cours...";
    try {
      const result = await invoke("ssh_connect", {
        profile: sshProfile,
        host: sshHost,
        user: sshUser,
        port: sshPort,
        password: sshPassword,
        keyPath: sshKeyPath
      });
      connectStatus = result;
    } catch (error) {
      connectStatus = `Erreur: ${error}`;
    }
  }

  async function scpUpload(event) {
    event.preventDefault();
    scpStatus = "Upload en cours...";
    try {
      const result = await invoke("scp_upload", {
        profile: scpProfile,
        host: scpHost,
        user: scpUser,
        port: scpPort,
        password: scpPassword,
        keyPath: scpKeyPath,
        localPath: scpLocalPath,
        remotePath: scpRemotePath
      });
      scpStatus = result;
    } catch (error) {
      scpStatus = `Erreur: ${error}`;
    }
  }

  async function scpDownload(event) {
    event.preventDefault();
    scpStatus = "Download en cours...";
    try {
      const result = await invoke("scp_download", {
        profile: scpProfile,
        host: scpHost,
        user: scpUser,
        port: scpPort,
        password: scpPassword,
        keyPath: scpKeyPath,
        localPath: scpLocalPath,
        remotePath: scpRemotePath
      });
      scpStatus = result;
    } catch (error) {
      scpStatus = `Erreur: ${error}`;
    }
  }

  async function loadProfiles() {
    try {
      profiles = await invoke("profile_list");
    } catch (error) {
      profileStatus = `Erreur: ${error}`;
    }
  }

  async function addProfile(event) {
    event.preventDefault();
    profileStatus = "Ajout du profil...";
    try {
      const result = await invoke("profile_add", {
        name: profileName,
        host: profileHost,
        user: profileUser,
        port: profilePort,
        password: profilePassword,
        keyPath: profileKeyPath
      });
      profileStatus = result;
      await loadProfiles();
      // Reset form
      profileName = "";
      profileHost = "";
      profileUser = "";
      profilePort = 22;
      profilePassword = "";
      profileKeyPath = "";
    } catch (error) {
      profileStatus = `Erreur: ${error}`;
    }
  }

  async function deleteProfile(name) {
    try {
      await invoke("profile_delete", { name });
      await loadProfiles();
      profileStatus = `Profil "${name}" supprimé`;
    } catch (error) {
      profileStatus = `Erreur: ${error}`;
    }
  }

  $effect(() => {
    if (currentTab === "profiles") {
      loadProfiles();
    }
  });
</script>

<main class="container">
  <h1>MySSH GUI</h1>

  <div class="tabs">
    <button
      class:active={currentTab === "connect"}
      onclick={() => currentTab = "connect"}>
      SSH Connect
    </button>
    <button
      class:active={currentTab === "scp"}
      onclick={() => currentTab = "scp"}>
      SCP Transfer
    </button>
    <button
      class:active={currentTab === "profiles"}
      onclick={() => currentTab = "profiles"}>
      Profils
    </button>
  </div>

  {#if currentTab === "connect"}
    <div class="tab-content">
      <h2>Connexion SSH</h2>
      <form onsubmit={connectSSH}>
        <div class="form-group">
          <label for="ssh-profile">Profil (optionnel)</label>
          <input id="ssh-profile" bind:value={sshProfile} placeholder="Nom du profil" />
        </div>

        <div class="divider">OU</div>

        <div class="form-group">
          <label for="ssh-host">Host *</label>
          <input id="ssh-host" bind:value={sshHost} placeholder="example.com" required={!sshProfile} />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="ssh-user">User *</label>
            <input id="ssh-user" bind:value={sshUser} placeholder="user" required={!sshProfile} />
          </div>

          <div class="form-group">
            <label for="ssh-port">Port</label>
            <input id="ssh-port" type="number" bind:value={sshPort} placeholder="22" />
          </div>
        </div>

        <div class="form-group">
          <label for="ssh-password">Password</label>
          <input id="ssh-password" type="password" bind:value={sshPassword} placeholder="Mot de passe" />
        </div>

        <div class="form-group">
          <label for="ssh-key">Clé privée SSH</label>
          <input id="ssh-key" bind:value={sshKeyPath} placeholder="/path/to/key" />
        </div>

        <button type="submit">Connecter</button>
        {#if connectStatus}
          <p class="status">{connectStatus}</p>
        {/if}
      </form>
    </div>
  {/if}

  {#if currentTab === "scp"}
    <div class="tab-content">
      <h2>Transfert SCP</h2>

      <div class="form-group">
        <label for="scp-profile">Profil (optionnel)</label>
        <input id="scp-profile" bind:value={scpProfile} placeholder="Nom du profil" />
      </div>

      <div class="divider">OU</div>

      <div class="form-group">
        <label for="scp-host">Host *</label>
        <input id="scp-host" bind:value={scpHost} placeholder="example.com" />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="scp-user">User *</label>
          <input id="scp-user" bind:value={scpUser} placeholder="user" />
        </div>

        <div class="form-group">
          <label for="scp-port">Port</label>
          <input id="scp-port" type="number" bind:value={scpPort} placeholder="22" />
        </div>
      </div>

      <div class="form-group">
        <label for="scp-password">Password</label>
        <input id="scp-password" type="password" bind:value={scpPassword} placeholder="Mot de passe" />
      </div>

      <div class="form-group">
        <label for="scp-key">Clé privée SSH</label>
        <input id="scp-key" bind:value={scpKeyPath} placeholder="/path/to/key" />
      </div>

      <hr />

      <div class="form-row">
        <div class="form-group">
          <label for="scp-local">Fichier local</label>
          <input id="scp-local" bind:value={scpLocalPath} placeholder="/local/path" />
        </div>

        <div class="form-group">
          <label for="scp-remote">Fichier distant</label>
          <input id="scp-remote" bind:value={scpRemotePath} placeholder="/remote/path" />
        </div>
      </div>

      <div class="button-group">
        <button onclick={scpUpload}>Upload →</button>
        <button onclick={scpDownload}>← Download</button>
      </div>

      {#if scpStatus}
        <p class="status">{scpStatus}</p>
      {/if}
    </div>
  {/if}

  {#if currentTab === "profiles"}
    <div class="tab-content">
      <h2>Gestion des profils</h2>

      <h3>Ajouter un profil</h3>
      <form onsubmit={addProfile}>
        <div class="form-group">
          <label for="profile-name">Nom *</label>
          <input id="profile-name" bind:value={profileName} placeholder="mon-serveur" required />
        </div>

        <div class="form-group">
          <label for="profile-host">Host *</label>
          <input id="profile-host" bind:value={profileHost} placeholder="example.com" required />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="profile-user">User *</label>
            <input id="profile-user" bind:value={profileUser} placeholder="user" required />
          </div>

          <div class="form-group">
            <label for="profile-port">Port</label>
            <input id="profile-port" type="number" bind:value={profilePort} placeholder="22" />
          </div>
        </div>

        <div class="form-group">
          <label for="profile-password">Password</label>
          <input id="profile-password" type="password" bind:value={profilePassword} placeholder="Mot de passe" />
        </div>

        <div class="form-group">
          <label for="profile-key">Clé privée SSH</label>
          <input id="profile-key" bind:value={profileKeyPath} placeholder="/path/to/key" />
        </div>

        <button type="submit">Ajouter</button>
      </form>

      {#if profileStatus}
        <p class="status">{profileStatus}</p>
      {/if}

      <hr />

      <h3>Profils enregistrés</h3>
      {#if profiles.length === 0}
        <p>Aucun profil enregistré</p>
      {:else}
        <div class="profile-list">
          {#each profiles as profile}
            <div class="profile-item">
              <div class="profile-info">
                <strong>{profile.name}</strong>
                <span>{profile.user}@{profile.host}:{profile.port}</span>
              </div>
              <button class="delete-btn" onclick={() => deleteProfile(profile.name)}>
                Supprimer
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</main>

<style>
  :root {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    font-size: 14px;
    line-height: 1.6;
    color: #333;
    background-color: #f5f5f5;
  }

  .container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  h1 {
    text-align: center;
    color: #2c3e50;
    margin-bottom: 30px;
  }

  h2 {
    color: #34495e;
    margin-bottom: 20px;
  }

  h3 {
    color: #34495e;
    margin-top: 20px;
    margin-bottom: 15px;
  }

  .tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
    border-bottom: 2px solid #ddd;
  }

  .tabs button {
    padding: 10px 20px;
    border: none;
    background: none;
    cursor: pointer;
    font-size: 16px;
    color: #7f8c8d;
    border-bottom: 3px solid transparent;
    transition: all 0.3s;
  }

  .tabs button:hover {
    color: #2c3e50;
  }

  .tabs button.active {
    color: #3498db;
    border-bottom-color: #3498db;
  }

  .tab-content {
    background: white;
    padding: 30px;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
  }

  label {
    display: block;
    margin-bottom: 5px;
    font-weight: 600;
    color: #555;
  }

  input {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
    transition: border-color 0.3s;
  }

  input:focus {
    outline: none;
    border-color: #3498db;
  }

  button {
    padding: 10px 20px;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 600;
    transition: background-color 0.3s;
  }

  button:hover {
    background-color: #2980b9;
  }

  .button-group {
    display: flex;
    gap: 10px;
  }

  .divider {
    text-align: center;
    margin: 20px 0;
    color: #999;
    font-size: 12px;
    text-transform: uppercase;
  }

  .status {
    margin-top: 15px;
    padding: 10px;
    background-color: #ecf0f1;
    border-radius: 4px;
    color: #2c3e50;
  }

  hr {
    margin: 30px 0;
    border: none;
    border-top: 1px solid #ddd;
  }

  .profile-list {
    margin-top: 20px;
  }

  .profile-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
    background-color: #f8f9fa;
    border-radius: 4px;
    margin-bottom: 10px;
  }

  .profile-info {
    display: flex;
    flex-direction: column;
  }

  .profile-info strong {
    color: #2c3e50;
    margin-bottom: 5px;
  }

  .profile-info span {
    color: #7f8c8d;
    font-size: 13px;
  }

  .delete-btn {
    background-color: #e74c3c;
    padding: 8px 16px;
  }

  .delete-btn:hover {
    background-color: #c0392b;
  }

  @media (prefers-color-scheme: dark) {
    :root {
      color: #ecf0f1;
      background-color: #2c3e50;
    }

    h1, h2, h3 {
      color: #ecf0f1;
    }

    .tab-content {
      background: #34495e;
    }

    input {
      background-color: #2c3e50;
      border-color: #465669;
      color: #ecf0f1;
    }

    .tabs button {
      color: #95a5a6;
    }

    .tabs button:hover {
      color: #ecf0f1;
    }

    .profile-item {
      background-color: #2c3e50;
    }

    .status {
      background-color: #2c3e50;
      color: #ecf0f1;
    }
  }
</style>
