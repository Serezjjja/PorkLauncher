import type { Translations } from "../types";

export const en: Translations = {
  common: {
    play: "PLAY",
    install: "INSTALL...",
    ready: "Ready",
    cancel: "Cancel",
    close: "Close",
    delete: "Delete",
    confirm: "Confirm",
    update: "Update",
    updateAvailable: "Update available",
    updating: "Updating",
    error: "Error",
    copy: "Copy",
    copied: "Copied!",
  },
  pages: {
    home: "Home",
    servers: "Servers",
    mods: "Mods",
  },
  profile: {
    username: "Username",
    version: "Version",
    noVersion: "No",
    releaseType: {
      preRelease: "Pre-Release",
      release: "Release",
    },
    loading: "Loading",
  },
  control: {
    status: {
      readyToPlay: "Ready to play",
    },
    updateAvailable: "Update available",
  },
  modals: {
    delete: {
      title: "Are you sure?",
      message: "Do you really want to delete the game?",
      warning:
        "This action will delete all game files without the possibility of recovery!",
      confirmButton: "Delete all",
      cancelButton: "Cancel",
    },
    error: {
      title: "An error occurred",
      technicalDetails: "Technical details",
      stackTrace: "Stack trace",
      suggestion: "Please report this issue if it persists.",
      copyError: "Copy error",
      copied: "Copied!",
      suggestions: {
        network: "Check your internet connection and try again.",
        filesystem:
          "Make sure you have enough disk space and the launcher has proper permissions.",
        validation: "Please check your input and try again.",
        game: "Try restarting the launcher or reinstalling the game.",
        default: "Please report this issue if it persists.",
      },
    },
    update: {
      title: "UPDATING LAUNCHER",
      message:
        "Downloading the latest version. PorkLand Launcher will restart automatically once finished.",
    },
    server: {
      copyIp: "Copy IP",
      copied: "Copied!",
      play: "Play",
    },
  },
  banners: {
    advertising: "Contact @hylauncher_bot for advertising",
    noServers: "No servers available",
    hynexus: {
      text: "HyNexus - this is Hytale as it should be. Economy, Clans, PVP, PVE, we're waiting for you!",
    },
    nctale: {
      text: "Join the community on Discord and Telegram — news, help and chat.",
    },
  },
  settings: {
    title: "SETTINGS",
    sections: {
      storage: "Storage",
      privacy: "Privacy",
      language: "Language",
    },
    storage: {
      logs: "Logs",
      logsDescription: "Browse or clean up your log files.",
      openLogs: "Open logs",
      deleteLogs: "Delete logs",
      clearCache: "Clear Cache/Game",
      clearCacheDescription: "Clean up PorkLand Launcher cache or game files. (will temporarily increase launch time)",
      deleteCache: "Delete Cache",
      deleteFiles: "Delete Files",
    },
    privacy: {
      analytics: "Analytics",
      analyticsDescription: "PorkLand Launcher collects analytics to improve the user experience.",
      discordRPC: "Discord RPC",
      discordRPCDescription: "Disabling this will cause 'PorkLand Launcher' to no longer show up as a game or app you are using on your Discord profile.",
    },
    language: {
      note: "Note:",
      translationNotice: "The app is not fully translated yet, so some content may remain in English for certain languages.",
    },
  },
  auth: {
    loginTitle: "Login",
    loginSubtitle: "Sign in with your account to play",
    emailLabel: "Email",
    emailPlaceholder: "your@email.com",
    passwordLabel: "Password",
    passwordPlaceholder: "••••••••",
    loginButton: "Sign In",
    loggingIn: "Signing in...",
    noAccount: "Don't have an account?",
    registerLink: "Register",
    logout: "Logout",
    welcome: "Welcome",
    noAccess: "No access to game",
    checkAccount: "Check your account on the website",
  },
};
