import React, { useState, useEffect, useLayoutEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  X,
  HardDrive,
  Shield,
  Languages,
  FolderOpen,
  FolderSearch,
  Trash2,
  Check,
} from "lucide-react";
import bgSettingsImage from "../assets/images/bg-settings.png";
import kweebecLogo from "../assets/images/kweeby.png";
import { useTranslation } from "../i18n";

export type SettingsSection = "storage" | "privacy" | "language";

interface SettingsOverlayProps {
  isOpen: boolean;
  onClose: () => void;
  launcherVersion: string;
}

export const SettingsOverlay: React.FC<SettingsOverlayProps> = ({
  isOpen,
  onClose,
  launcherVersion,
}) => {
  const { t } = useTranslation();
  const [overlayEnterReady, setOverlayEnterReady] = useState(false);
  const [settingsSection, setSettingsSection] =
    useState<SettingsSection>("storage");

  useEffect(() => {
    if (!isOpen) {
      setOverlayEnterReady(false);
    } else {
      setSettingsSection("storage");
    }
  }, [isOpen]);

  useLayoutEffect(() => {
    if (!isOpen) return;
    setOverlayEnterReady(false);
    const raf = requestAnimationFrame(() => setOverlayEnterReady(true));
    return () => cancelAnimationFrame(raf);
  }, [isOpen]);

  // Close on ESC
  useEffect(() => {
    if (!isOpen) return;
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [isOpen, onClose]);

  return (
    <AnimatePresence mode="wait">
      {isOpen && (
        <motion.div
          key="settings-overlay"
          initial={{ opacity: 1 }}
          animate={{ opacity: 1 }}
          exit={{
            opacity: 0,
            scale: 0.98,
            transition: { duration: 0.28, ease: [0.16, 1, 0.3, 1] },
          }}
          className="absolute inset-0 z-[100] flex items-center justify-center origin-center"
          tabIndex={0}
          onKeyDown={(e) => e.key === "Escape" && onClose()}
        >
          {/* Background blur */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{
              opacity: overlayEnterReady ? 1 : 0,
              transition: { duration: 0.4, ease: [0.22, 1, 0.36, 1] },
            }}
            className="absolute inset-0 bg-black/25 backdrop-blur-[14px] cursor-default"
            onClick={onClose}
            aria-hidden
          />
          {/* Main window */}
          <motion.div
            initial={{ opacity: 0, scale: 0.97 }}
            animate={{
              opacity: overlayEnterReady ? 1 : 0,
              scale: overlayEnterReady ? 1 : 0.97,
              transition: {
                duration: 0.22,
                delay: 0,
                ease: [0.16, 1, 0.3, 1],
              },
            }}
            className="relative z-10 w-[900px] h-[500px] rounded-[20px] border border-[#5D8B57]/[0.3] bg-cover bg-center bg-no-repeat overflow-hidden shadow-[0_24px_64px_rgba(0,0,0,0.5),0_0_40px_rgba(93,139,87,0.12)]"
            style={{ backgroundImage: `url(${bgSettingsImage})` }}
            onClick={(e) => e.stopPropagation()}
            role="dialog"
            aria-label="Settings"
          >
            {/* Dark overlay for readability */}
            <div
              className="absolute inset-0 bg-[#1C2A1F]/[0.85] rounded-[20px]"
              aria-hidden
            />

            {/* Header */}
            <div className="absolute left-[30px] top-[30px] z-10 flex items-center gap-[12px]">
              <span className="text-[20px] font-[Unbounded] font-[500] uppercase tracking-wide text-white/90">
                {t.settings?.title || "SETTINGS"}
              </span>
              <span
                className="w-[1px] h-[20px] bg-[#7C7C7C]/[0.10]"
                aria-hidden
              />
              <span className="text-[14px] text-white/25 font-[Mazzard]">
                PorkLand Launcher v{launcherVersion}
              </span>
            </div>

            {/* Horizontal divider */}
            <div
              className="absolute left-0 right-0 top-[80px] z-10 h-[1px] bg-gradient-to-r from-transparent via-[#5D8B57]/30 to-transparent"
              aria-hidden
            />

            {/* Vertical divider */}
            <div
              className="absolute left-[176px] top-[81px] bottom-0 z-10 w-[1px] bg-gradient-to-b from-[#5D8B57]/30 via-[#5D8B57]/20 to-transparent"
              aria-hidden
            />

            {/* Sidebar */}
            <div className="absolute left-[30px] top-[111px] z-10 flex flex-col gap-[12px]">
              <SidebarButton
                icon={<HardDrive size={18} strokeWidth={2} />}
                label={t.settings?.sections?.storage || "Storage"}
                isActive={settingsSection === "storage"}
                onClick={() => setSettingsSection("storage")}
              />
              <SidebarButton
                icon={<Shield size={18} strokeWidth={2} />}
                label={t.settings?.sections?.privacy || "Privacy"}
                isActive={settingsSection === "privacy"}
                onClick={() => setSettingsSection("privacy")}
              />
              <SidebarButton
                icon={<Languages size={18} strokeWidth={2} />}
                label={t.settings?.sections?.language || "Language"}
                isActive={settingsSection === "language"}
                onClick={() => setSettingsSection("language")}
              />
            </div>

            {/* Content area */}
            <div className="absolute left-[206px] right-[30px] top-[111px] bottom-[30px] z-10 overflow-y-auto">
              {settingsSection === "storage" && <StorageSection />}
              {settingsSection === "privacy" && <PrivacySection />}
              {settingsSection === "language" && <LanguageSection />}
            </div>

            {/* Close button */}
            <button
              type="button"
              onClick={onClose}
              className="absolute right-[30px] top-[30px] z-10 flex items-center justify-center text-[#999999] transition-colors hover:text-white cursor-pointer"
              aria-label="Close settings"
            >
              <X size={18} strokeWidth={2} />
            </button>

            {/* Footer text */}
            <div
              className="absolute left-[30px] bottom-[30px] z-10 text-[12px] font-[Unbounded] text-white/25"
              aria-hidden
            >
              {"PorkLand Launcher <3"}
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

interface SidebarButtonProps {
  icon: React.ReactNode;
  label: string;
  isActive: boolean;
  onClick: () => void;
}

const SidebarButton: React.FC<SidebarButtonProps> = ({
  icon,
  label,
  isActive,
  onClick,
}) => (
  <button
    type="button"
    onClick={onClick}
    className={`flex items-center gap-2 px-3 py-2 rounded-xl text-left transition-all duration-300 cursor-pointer ${
      isActive 
        ? "text-[#7AB872] bg-[#5D8B57]/[0.15] shadow-[0_0_15px_rgba(93,139,87,0.2)] drop-shadow-[0_0_6px_rgba(122,184,114,0.5)]" 
        : "text-white/60 hover:text-[#7AB872] hover:bg-[#5D8B57]/[0.08]"
    }`}
    aria-label={label}
  >
    {icon}
    <span className="text-[16px] font-[Mazzard]">{label}</span>
  </button>
);

const StorageSection: React.FC = () => {
  const { t } = useTranslation();
  const [gameDir, setGameDir] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);

  // Load game directory on mount
  useEffect(() => {
    const loadGameDir = async () => {
      try {
        // @ts-ignore - Wails bindings
        const result = await window.go?.app?.App?.GetGameDirectory();
        if (result) {
          setGameDir(result);
        }
      } catch (err) {
        console.error("Failed to load game directory:", err);
      } finally {
        setIsLoading(false);
      }
    };
    loadGameDir();
  }, []);

  const handleBrowse = async () => {
    try {
      // @ts-ignore - Wails bindings
      const selected = await window.go?.app?.App?.BrowseGameDirectory();
      if (selected) {
        // @ts-ignore - Wails bindings
        await window.go?.app?.App?.SetGameDirectory(selected);
        setGameDir(selected);
      }
    } catch (err) {
      console.error("Failed to browse game directory:", err);
    }
  };

  const handleOpenLogs = async () => {
    try {
      // @ts-ignore - Wails bindings
      await window.go?.app?.App?.OpenLogsFolder();
    } catch (err) {
      console.error("Failed to open logs folder:", err);
    }
  };

  const handleDeleteLogs = async () => {
    try {
      // @ts-ignore - Wails bindings
      await window.go?.app?.App?.DeleteLogs();
      console.log("[StorageSection] Logs deleted successfully");
    } catch (err) {
      console.error("Failed to delete logs:", err);
    }
  };

  const handleDeleteCache = async () => {
    try {
      // @ts-ignore - Wails bindings
      await window.go?.app?.App?.DeleteCache();
      console.log("[StorageSection] Cache deleted successfully");
    } catch (err) {
      console.error("Failed to delete cache:", err);
    }
  };

  const handleDeleteFiles = async () => {
    try {
      // @ts-ignore - Wails bindings
      await window.go?.app?.App?.DeleteFiles();
      console.log("[StorageSection] Game files deleted successfully");
    } catch (err) {
      console.error("Failed to delete game files:", err);
    }
  };

  return (
    <div className="flex flex-col gap-[24px] text-white/90 font-[Mazzard]">
      {/* <section>
        <h3 className="text-[16px] font-[Unbounded] font-[500] text-white mb-[6px]">
          Game directory
        </h3>
        <p className="text-[14px] font-[Mazzard] text-white/50 mb-[6px]">
          The directory where the game stores all of its files. Changes will be
          applied after restarting the launcher.
        </p>
        <div className="relative w-full">
          <input
            type="text"
            readOnly
            className="w-full h-[46px] pl-4 pr-10 rounded-[14px] bg-[#090909]/[0.55] border border-[#7C7C7C]/[0.10] text-[14px] text-[#CCD9E0]/[0.9] font-[Mazzard]"
            value={isLoading ? "Loading..." : gameDir}
          />
          <button
            type="button"
            onClick={handleBrowse}
            className="absolute right-[16px] top-1/2 -translate-y-1/2 flex items-center justify-center w-8 h-8 text-[#CCD9E0]/[0.9] hover:opacity-80 transition-opacity cursor-pointer"
            aria-label="Browse"
          >
            <FolderSearch size={18} />
          </button>
        </div>
      </section> */}
      <section>
        <h3 className="text-[16px] font-[Unbounded] font-[500] text-white mb-[6px]">
          {t.settings?.storage?.logs || "Logs"}
        </h3>
        <p className="text-[14px] font-[Mazzard] text-white/50 mb-[6px]">
          {t.settings?.storage?.logsDescription || "Browse or clean up your log files."}
        </p>
        <div className="flex items-center gap-[10px]">
          <button
            type="button"
            onClick={handleOpenLogs}
            className="flex items-center justify-center gap-[16px] w-[130px] h-[46px] rounded-[14px] bg-[#364E3A]/[0.65] border border-[#5D8B57]/[0.3] font-[Mazzard] text-white/90 text-[14px] hover:bg-[#5D8B57]/[0.15] hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_15px_rgba(93,139,87,0.2)] transition-all duration-300 cursor-pointer"
          >
            {t.settings?.storage?.openLogs || "Open logs"} <FolderOpen size={16} />
          </button>
          <button
            type="button"
            onClick={handleDeleteLogs}
            className="flex items-center justify-center gap-[16px] w-[136px] h-[46px] rounded-[14px] bg-[#364E3A]/[0.65] border border-[#5D8B57]/[0.3] font-[Mazzard] text-white/90 text-[14px] hover:bg-[#5D8B57]/[0.15] hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_15px_rgba(93,139,87,0.2)] transition-all duration-300 cursor-pointer"
          >
            {t.settings?.storage?.deleteLogs || "Delete logs"} <Trash2 size={16} />
          </button>
        </div>
      </section>
      <section>
        <h3 className="text-[16px] font-[Unbounded] font-[500] text-white mb-[6px]">
          {t.settings?.storage?.clearCache || "Clear Cache/Game"}
        </h3>
        <p className="text-[14px] font-[Mazzard] text-white/50 mb-[6px]">
          {t.settings?.storage?.clearCacheDescription || "Clean up PorkLand Launcher cache or game files. (will temporarily increase launch time)"}
        </p>
        <div className="flex items-center gap-[10px]">
          <button
            type="button"
            onClick={handleDeleteCache}
            className="flex items-center justify-center gap-[16px] w-[154px] h-[46px] rounded-[14px] bg-[#364E3A]/[0.65] border border-[#5D8B57]/[0.3] font-[Mazzard] text-white/90 text-[14px] hover:bg-[#5D8B57]/[0.15] hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_15px_rgba(93,139,87,0.2)] transition-all duration-300 cursor-pointer"
          >
            {t.settings?.storage?.deleteCache || "Delete Cache"} <Trash2 size={16} />
          </button>
          <button
            type="button"
            onClick={handleDeleteFiles}
            className="flex items-center justify-center gap-[16px] w-[140px] h-[46px] rounded-[14px] bg-[#2A1A1A]/[0.75] border border-[#8B0000]/[0.3] font-[Mazzard] text-white text-[14px] hover:bg-[#3D1F1F]/[0.85] hover:border-[#A00000]/[0.5] hover:shadow-[0_0_15px_rgba(139,0,0,0.25)] transition-all duration-300 cursor-pointer"
          >
            {t.settings?.storage?.deleteFiles || "Delete Files"} <Trash2 size={16} />
          </button>
        </div>
      </section>
    </div>
  );
};

const PrivacySection: React.FC = () => {
  const { t } = useTranslation();
  const [analytics, setAnalytics] = useState(true);
  const [discordRPC, setDiscordRPC] = useState(true);
  const [isLoading, setIsLoading] = useState(true);

  // Load Discord RPC setting on mount
  useEffect(() => {
    const loadSetting = async () => {
      try {
        // @ts-ignore - Wails bindings
        const result = await window.go?.app?.App?.GetDiscordRPC();
        if (typeof result === "boolean") {
          setDiscordRPC(result);
        }
      } catch (err) {
        console.error("Failed to load Discord RPC setting:", err);
      } finally {
        setIsLoading(false);
      }
    };
    loadSetting();
  }, []);

  const handleDiscordRPCChange = async (enabled: boolean) => {
    try {
      // @ts-ignore - Wails bindings
      await window.go?.app?.App?.SetDiscordRPC(enabled);
      setDiscordRPC(enabled);
    } catch (err) {
      console.error("Failed to save Discord RPC setting:", err);
    }
  };

  return (
    <div className="flex flex-col gap-[32px] text-white/90 font-[Mazzard]">
      {/* Analytics */}
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-[16px] font-[Unbounded] font-[500] text-white mb-[8px]">
            {t.settings?.privacy?.analytics || "Analytics"}
          </h3>
          <p className="text-[14px] font-[Mazzard] text-white/50">
            {t.settings?.privacy?.analyticsDescription || "PorkLand Launcher collects analytics to improve the user experience."}
          </p>
        </div>
        <ToggleSwitch checked={analytics} onChange={setAnalytics} />
      </div>

      {/* Discord RPC */}
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-[16px] font-[Unbounded] font-[500] text-white mb-[8px]">
            {t.settings?.privacy?.discordRPC || "Discord RPC"}
          </h3>
          <p className="text-[14px] font-[Mazzard] text-white/50 max-w-[400px]">
            {t.settings?.privacy?.discordRPCDescription || "Disabling this will cause 'PorkLand Launcher' to no longer show up as a game or app you are using on your Discord profile."}
          </p>
        </div>
        <ToggleSwitch
          checked={discordRPC}
          onChange={handleDiscordRPCChange}
          disabled={isLoading}
        />
      </div>
    </div>
  );
};

interface ToggleSwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}

const ToggleSwitch: React.FC<ToggleSwitchProps> = ({
  checked,
  onChange,
  disabled,
}) => {
  return (
    <button
      type="button"
      onClick={() => !disabled && onChange(!checked)}
      disabled={disabled}
      className={`
        w-[48px] h-[48px] rounded-[14px] border
        flex items-center justify-center
        transition-all duration-200
        ${disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}
        ${
          checked
            ? "bg-[#364E3A]/[0.65] border-[#5D8B57]/[0.4] shadow-[0_0_15px_rgba(93,139,87,0.2)]"
            : "bg-[#364E3A]/[0.35] border-[#5D8B57]/[0.15]"
        }
      `}
      aria-label={checked ? "Enabled" : "Disabled"}
    >
      {checked && (
        <svg
          width="20"
          height="20"
          viewBox="0 0 20 20"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M4 10L8 14L16 6"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      )}
    </button>
  );
};

const LanguageSection: React.FC = () => {
  const { language, setLanguage, t } = useTranslation();

  const languages = [
    { code: "ru", name: "Русский" },
    { code: "en", name: "English" },
  ] as const;

  return (
    <div className="flex flex-col gap-[16px] text-white/90 font-[Mazzard]">
      {/* Notice block */}
      <div
        className="
        w-full h-[80px]
        bg-[#364E3A]/[0.65] backdrop-blur-[24px]
        border border-[#5D8B57]/[0.3]
        rounded-[14px]
        flex items-center
        px-[10px] gap-[10px]
        shadow-[0_8px_32px_rgba(0,0,0,0.3),0_0_20px_rgba(93,139,87,0.1)]
      "
      >
        <img
          src={kweebecLogo}
          alt="Kweebec"
          className="w-[60px] h-[60px] rounded-[8px]"
        />
        <div className="flex flex-col gap-[8px]">
          <span className="text-[14px] font-[Unbounded] font-[500] text-[#7AB872] tracking-[-0.03em] drop-shadow-[0_0_6px_rgba(122,184,114,0.4)]">
            {t.settings?.language?.note || "Note:"}
          </span>
          <span className="text-[14px] font-[Mazzard] font-[500] text-white/50 tracking-[-0.03em] leading-[110%]">
            {t.settings?.language?.translationNotice ||
              "The app is not fully translated yet, so some content may remain in English for certain languages."}
          </span>
        </div>
      </div>

      {/* Language list */}
      <div className="flex flex-col gap-[8px]">
        {languages.map((lang) => (
          <button
            key={lang.code}
            onClick={() => setLanguage(lang.code)}
            className="
              w-full h-[48px]
              bg-[#364E3A]/[0.65] backdrop-blur-[24px]
              border border-[#5D8B57]/[0.25]
              rounded-[14px]
              flex items-center justify-between
              px-[16px]
              cursor-pointer
              hover:bg-[#5D8B57]/[0.12]
              hover:border-[#5D8B57]/[0.45]
              hover:shadow-[0_0_15px_rgba(93,139,87,0.15)]
              transition-all duration-300
            "
          >
            <span className="text-[16px] font-[Mazzard] font-[600] text-white/90 tracking-[-0.03em]">
              {lang.name}
            </span>
            {language === lang.code && (
              <Check
                size={16}
                className="text-[#7AB872] drop-shadow-[0_0_6px_rgba(122,184,114,0.6)]"
                strokeWidth={2}
              />
            )}
          </button>
        ))}
      </div>
    </div>
  );
};

export default SettingsOverlay;
