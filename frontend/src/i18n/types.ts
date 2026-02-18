/**
 * Type definitions for i18n system
 * Keeps all translation logic separate from UI
 */

export type SupportedLanguage = "en" | "ru";

export interface Translations {
  common: {
    play: string;
    install: string;
    ready: string;
    cancel: string;
    close: string;
    delete: string;
    confirm: string;
    update: string;
    updateAvailable: string;
    updating: string;
    error: string;
    copy: string;
    copied: string;
  };
  pages: {
    home: string;
    servers: string;
    mods: string;
  };
  profile: {
    username: string;
    version: string;
    noVersion: string;
    releaseType: {
      preRelease: string;
      release: string;
    };
    loading: string;
  };
  control: {
    status: {
      readyToPlay: string;
    };
    updateAvailable: string;
  };
  modals: {
    delete: {
      title: string;
      message: string;
      warning: string;
      confirmButton: string;
      cancelButton: string;
    };
    error: {
      title: string;
      technicalDetails: string;
      stackTrace: string;
      suggestion: string;
      copyError: string;
      copied: string;
      suggestions: {
        network: string;
        filesystem: string;
        validation: string;
        game: string;
        default: string;
      };
    };
    update: {
      title: string;
      message: string;
    };
    server: {
      copyIp: string;
      copied: string;
      play: string;
    };
  };
  banners: {
    advertising: string;
    noServers: string;
    hynexus: {
      text: string;
    };
    nctale: {
      text: string;
    };
  };
  settings: {
    title: string;
    sections: {
      storage: string;
      privacy: string;
      language: string;
    };
    storage: {
      logs: string;
      logsDescription: string;
      openLogs: string;
      deleteLogs: string;
      clearCache: string;
      clearCacheDescription: string;
      deleteCache: string;
      deleteFiles: string;
    };
    privacy: {
      analytics: string;
      analyticsDescription: string;
      discordRPC: string;
      discordRPCDescription: string;
    };
    language: {
      note: string;
      translationNotice: string;
    };
  };
  auth: {
    loginTitle: string;
    loginSubtitle: string;
    emailLabel: string;
    emailPlaceholder: string;
    passwordLabel: string;
    passwordPlaceholder: string;
    loginButton: string;
    loggingIn: string;
    noAccount: string;
    registerLink: string;
    logout: string;
    welcome: string;
    noAccess: string;
    checkAccount: string;
  };
}

export interface I18nContextValue {
  language: SupportedLanguage;
  setLanguage: (lang: SupportedLanguage) => void;
  t: Translations;
}
