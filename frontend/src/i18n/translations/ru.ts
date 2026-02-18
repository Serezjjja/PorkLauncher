import type { Translations } from "../types";

export const ru: Translations = {
  common: {
    play: "ИГРАТЬ",
    install: "УСТАНОВКА...",
    ready: "Готово",
    cancel: "Отмена",
    close: "Закрыть",
    delete: "Удалить",
    confirm: "Подтвердить",
    update: "Обновить",
    updateAvailable: "Доступно обновление",
    updating: "Обновление",
    error: "Ошибка",
    copy: "Копировать",
    copied: "Скопировано!",
  },
  pages: {
    home: "Главная",
    servers: "Серверы",
    mods: "Моды",
  },
  profile: {
    username: "Имя пользователя",
    version: "Версия",
    noVersion: "Нет",
    releaseType: {
      preRelease: "Pre-Release",
      release: "Release",
    },
    loading: "Загрузка",
  },
  control: {
    status: {
      readyToPlay: "Готов к игре",
    },
    updateAvailable: "Доступно обновление",
  },
  modals: {
    delete: {
      title: "Вы уверены?",
      message: "Вы точно хотите удалить игру?",
      warning:
        "Это действие удалит все файлы игры без возможности восстановления!",
      confirmButton: "Удалить всё",
      cancelButton: "Отмена",
    },
    error: {
      title: "Произошла ошибка",
      technicalDetails: "Технические детали",
      stackTrace: "Трассировка стека",
      suggestion: "Пожалуйста, сообщите об этой проблеме, если она сохраняется.",
      copyError: "Скопировать ошибку",
      copied: "Скопировано!",
      suggestions: {
        network: "Проверьте подключение к интернету и попробуйте снова.",
        filesystem:
          "Убедитесь, что у вас достаточно места на диске и у лаунчера есть необходимые права доступа.",
        validation:
          "Пожалуйста, проверьте введенные данные и попробуйте снова.",
        game: "Попробуйте перезапустить лаунчер или переустановить игру.",
        default: "Пожалуйста, сообщите об этой проблеме, если она сохраняется.",
      },
    },
    update: {
      title: "ОБНОВЛЕНИЕ ЛАУНЧЕРА",
      message:
        "Загрузка последней версии. PorkLand Launcher автоматически перезапустится после завершения.",
    },
    server: {
      copyIp: "Копировать IP",
      copied: "Скопировано!",
      play: "Играть",
    },
  },
  banners: {
    advertising: "По поводу рекламы пишите нашему боту @hylauncher_bot",
    noServers: "Нет доступных серверов",
    hynexus: {
      text: "HyNexus - это Hytale, каким он должен быть. Экономика, Кланы, PVP, PVE, ждем тебя!",
    },
    nctale: {
      text: "NCTale — королевство в hytale! PvP-битвы, войны за территорию, варварство, экономика.",
    },
  },
  settings: {
    title: "НАСТРОЙКИ",
    sections: {
      storage: "Хранилище",
      privacy: "Приватность",
      language: "Язык",
    },
    storage: {
      logs: "Логи",
      logsDescription: "Просмотрите или очистите файлы логов.",
      openLogs: "Открыть логи",
      deleteLogs: "Удалить логи",
      clearCache: "Очистить кэш/игру",
      clearCacheDescription: "Очистите кэш PorkLand Launcher или файлы игры. (временно увеличит время запуска)",
      deleteCache: "Удалить кэш",
      deleteFiles: "Удалить файлы",
    },
    privacy: {
      analytics: "Аналитика",
      analyticsDescription: "PorkLand Launcher собирает аналитику для улучшения пользовательского опыта.",
      discordRPC: "Discord RPC",
      discordRPCDescription: "Отключение этой функции уберёт 'PorkLand Launcher' из статуса активности в вашем профиле Discord.",
    },
    language: {
      note: "Примечание:",
      translationNotice: "Приложение еще не полностью переведено, поэтому для некоторых языков часть контента может оставаться на английском языке.",
    },
  },
  auth: {
    loginTitle: "Вход",
    loginSubtitle: "Войдите в аккаунт, чтобы играть",
    emailLabel: "Email",
    emailPlaceholder: "your@email.com",
    passwordLabel: "Пароль",
    passwordPlaceholder: "••••••••",
    loginButton: "Войти",
    loggingIn: "Вход...",
    noAccount: "Нет аккаунта?",
    registerLink: "Регистрация",
    logout: "Выйти",
    welcome: "Добро пожаловать",
    noAccess: "Нет доступа к игре",
    checkAccount: "Проверьте свой аккаунт на сайте",
  },
};
