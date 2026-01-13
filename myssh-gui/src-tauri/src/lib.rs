use std::process::Command;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct Profile {
    name: String,
    host: String,
    port: i32,
    user: String,
}

// Helper function to get the myssh binary path
fn get_myssh_path() -> String {
    // In development, use the binary in the parent directory
    // In production, the binary should be bundled with the app
    std::env::var("MYSSH_PATH")
        .unwrap_or_else(|_| "../myssh".to_string())
}

#[tauri::command]
fn ssh_connect(
    profile: String,
    host: String,
    user: String,
    port: i32,
    password: String,
    key_path: String,
) -> Result<String, String> {
    let myssh_path = get_myssh_path();
    let mut cmd = Command::new(&myssh_path);
    cmd.arg("connect");

    if !profile.is_empty() {
        cmd.arg("--profile").arg(&profile);
    } else {
        if !host.is_empty() {
            cmd.arg("--host").arg(&host);
        }
        if !user.is_empty() {
            cmd.arg("--user").arg(&user);
        }
        if port > 0 {
            cmd.arg("--port").arg(port.to_string());
        }
        if !password.is_empty() {
            cmd.arg("--password").arg(&password);
        }
        if !key_path.is_empty() {
            cmd.arg("--key").arg(&key_path);
        }
    }

    match cmd.output() {
        Ok(output) => {
            if output.status.success() {
                Ok(String::from_utf8_lossy(&output.stdout).to_string())
            } else {
                Err(String::from_utf8_lossy(&output.stderr).to_string())
            }
        }
        Err(e) => Err(format!("Failed to execute myssh: {}", e)),
    }
}

#[tauri::command]
fn scp_upload(
    profile: String,
    host: String,
    user: String,
    port: i32,
    password: String,
    key_path: String,
    local_path: String,
    remote_path: String,
) -> Result<String, String> {
    let myssh_path = get_myssh_path();
    let mut cmd = Command::new(&myssh_path);
    cmd.arg("scp").arg("upload");

    if !profile.is_empty() {
        cmd.arg("--profile").arg(&profile);
    } else {
        if !host.is_empty() {
            cmd.arg("--host").arg(&host);
        }
        if !user.is_empty() {
            cmd.arg("--user").arg(&user);
        }
        if port > 0 {
            cmd.arg("--port").arg(port.to_string());
        }
        if !password.is_empty() {
            cmd.arg("--password").arg(&password);
        }
        if !key_path.is_empty() {
            cmd.arg("--key").arg(&key_path);
        }
    }

    cmd.arg("--local").arg(&local_path);
    cmd.arg("--remote").arg(&remote_path);

    match cmd.output() {
        Ok(output) => {
            if output.status.success() {
                Ok("Upload réussi !".to_string())
            } else {
                Err(String::from_utf8_lossy(&output.stderr).to_string())
            }
        }
        Err(e) => Err(format!("Failed to execute myssh: {}", e)),
    }
}

#[tauri::command]
fn scp_download(
    profile: String,
    host: String,
    user: String,
    port: i32,
    password: String,
    key_path: String,
    local_path: String,
    remote_path: String,
) -> Result<String, String> {
    let myssh_path = get_myssh_path();
    let mut cmd = Command::new(&myssh_path);
    cmd.arg("scp").arg("download");

    if !profile.is_empty() {
        cmd.arg("--profile").arg(&profile);
    } else {
        if !host.is_empty() {
            cmd.arg("--host").arg(&host);
        }
        if !user.is_empty() {
            cmd.arg("--user").arg(&user);
        }
        if port > 0 {
            cmd.arg("--port").arg(port.to_string());
        }
        if !password.is_empty() {
            cmd.arg("--password").arg(&password);
        }
        if !key_path.is_empty() {
            cmd.arg("--key").arg(&key_path);
        }
    }

    cmd.arg("--local").arg(&local_path);
    cmd.arg("--remote").arg(&remote_path);

    match cmd.output() {
        Ok(output) => {
            if output.status.success() {
                Ok("Download réussi !".to_string())
            } else {
                Err(String::from_utf8_lossy(&output.stderr).to_string())
            }
        }
        Err(e) => Err(format!("Failed to execute myssh: {}", e)),
    }
}

#[tauri::command]
fn profile_list() -> Result<Vec<Profile>, String> {
    let myssh_path = get_myssh_path();
    let output = Command::new(&myssh_path)
        .arg("profile")
        .arg("list")
        .output()
        .map_err(|e| format!("Failed to execute myssh: {}", e))?;

    if !output.status.success() {
        return Err(String::from_utf8_lossy(&output.stderr).to_string());
    }

    let output_str = String::from_utf8_lossy(&output.stdout);
    let mut profiles = Vec::new();

    // Parse the output (format: NAME HOST USER PORT)
    for line in output_str.lines().skip(1) {
        // Skip header
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() >= 4 {
            profiles.push(Profile {
                name: parts[0].to_string(),
                host: parts[1].to_string(),
                user: parts[2].to_string(),
                port: parts[3].parse().unwrap_or(22),
            });
        }
    }

    Ok(profiles)
}

#[tauri::command]
fn profile_add(
    name: String,
    host: String,
    user: String,
    port: i32,
    password: String,
    key_path: String,
) -> Result<String, String> {
    let myssh_path = get_myssh_path();
    let mut cmd = Command::new(&myssh_path);
    cmd.arg("profile").arg("add").arg(&name);
    cmd.arg("--host").arg(&host);
    cmd.arg("--user").arg(&user);
    cmd.arg("--port").arg(port.to_string());

    if !password.is_empty() {
        cmd.arg("--password").arg(&password);
    }
    if !key_path.is_empty() {
        cmd.arg("--key").arg(&key_path);
    }

    match cmd.output() {
        Ok(output) => {
            if output.status.success() {
                Ok(format!("Profil '{}' créé avec succès", name))
            } else {
                Err(String::from_utf8_lossy(&output.stderr).to_string())
            }
        }
        Err(e) => Err(format!("Failed to execute myssh: {}", e)),
    }
}

#[tauri::command]
fn profile_delete(name: String) -> Result<String, String> {
    let myssh_path = get_myssh_path();
    let output = Command::new(&myssh_path)
        .arg("profile")
        .arg("delete")
        .arg(&name)
        .output()
        .map_err(|e| format!("Failed to execute myssh: {}", e))?;

    if output.status.success() {
        Ok(format!("Profil '{}' supprimé", name))
    } else {
        Err(String::from_utf8_lossy(&output.stderr).to_string())
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![
            ssh_connect,
            scp_upload,
            scp_download,
            profile_list,
            profile_add,
            profile_delete
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
