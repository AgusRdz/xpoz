#Requires -Version 5.1
$ErrorActionPreference = "Stop"

$Repo   = "AgusRdz/xpoz"
$Binary = "xpoz"

# ── architecture ──────────────────────────────────────────────────────────────

$Arch = if (
  [System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture -eq
  [System.Runtime.InteropServices.Architecture]::Arm64
) { "arm64" } else { "amd64" }

# ── install directory ─────────────────────────────────────────────────────────

$InstallDir = if ($env:XPOZ_INSTALL_DIR) {
  $env:XPOZ_INSTALL_DIR
} else {
  Join-Path $env:LOCALAPPDATA "Programs\xpoz"
}

if (-not (Test-Path $InstallDir)) {
  New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# ── latest version ────────────────────────────────────────────────────────────

$Version = if ($env:XPOZ_VERSION) {
  $env:XPOZ_VERSION
} else {
  $release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
  $release.tag_name
}

if (-not $Version) {
  Write-Error "Could not determine latest version. Set XPOZ_VERSION to override."
  exit 1
}

# ── download and extract ──────────────────────────────────────────────────────

$Archive  = "$Binary-windows-$Arch.zip"
$Url      = "https://github.com/$Repo/releases/download/$Version/$Archive"
$Dest     = Join-Path $InstallDir "$Binary.exe"
$TmpZip   = Join-Path $env:TEMP "$Archive"
$TmpDir   = Join-Path $env:TEMP "xpoz-install-$Version"

Write-Host "  Downloading xpoz $Version (windows/$Arch)..."
Invoke-WebRequest -Uri $Url -OutFile $TmpZip -UseBasicParsing
Expand-Archive -Path $TmpZip -DestinationPath $TmpDir -Force
Move-Item (Join-Path $TmpDir "$Binary.exe") $Dest -Force
Remove-Item $TmpZip, $TmpDir -Recurse -Force -ErrorAction SilentlyContinue

# ── PATH setup ────────────────────────────────────────────────────────────────

$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
  [Environment]::SetEnvironmentVariable("PATH", "$InstallDir;$UserPath", "User")

  # Broadcast WM_SETTINGCHANGE so open terminals pick up the new PATH without restart
  try {
    $Signature = @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
  IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
  uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
    $Type   = Add-Type -MemberDefinition $Signature -Name "Win32SendMsg" -Namespace "Win32" -PassThru
    $Result = [UIntPtr]::Zero
    $Type::SendMessageTimeout([IntPtr]0xffff, 0x001a, [UIntPtr]::Zero, "Environment", 2, 5000, [ref]$Result) | Out-Null
  } catch {
    # Non-fatal — PATH will be active in new terminals
  }

  $env:PATH = "$InstallDir;$env:PATH"
  Write-Host "  Added $InstallDir to PATH"
}

Write-Host "✓ xpoz $Version installed to $Dest"
Write-Host ""
Write-Host "  Run 'xpoz setup' to get started."
