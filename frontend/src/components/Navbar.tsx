// src/components/Navbar.tsx
import React, { useContext } from "react";
import { getPages } from "../config/pages";
import { useTranslation } from "../i18n";
import telegramIcon from "../assets/images/telegram.svg";
import discordIcon from "../assets/images/discord.svg";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { OpenFolder } from "../../wailsjs/go/app/App";
import { Activity, Bolt, FolderOpen } from "lucide-react";
import { SettingsOverlayContext } from "../context/SettingsOverlayContext";

interface NavbarProps {
  activeTab: string;
  onTabChange: (tab: string) => void;
  onDiagnosticsClick?: () => void;
  onSettingsClick?: () => void;
}

function Navbar({ activeTab, onTabChange, onDiagnosticsClick, onSettingsClick }: NavbarProps) {
  const { t } = useTranslation();
  const pages = getPages(t);
  const isSettingsOpen = useContext(SettingsOverlayContext);

  const openTelegram = () => {
    try {
      BrowserOpenURL("https://t.me/porklandmc");
    } catch {
      window.open("https://t.me/porklandmc", "_blank");
    }
  };

  const openDiscord = () => {
    try {
      BrowserOpenURL("https://dsc.gg/porklandmc");
    } catch {
      window.open("https://dsc.gg/porklandmc", "_blank");
    }
  };

  const openWebsite = () => {
    try {
      BrowserOpenURL("https://porkland.net");
    } catch {
      window.open("https://porkland.net", "_blank");
    }
  };

  return (
    <div
      className="
        absolute left-[24px] top-1/2 -translate-y-1/2
        w-[56px] h-[340px]
        bg-[#364E3A]/[0.75]
        backdrop-blur-[24px]
        rounded-[20px]
        border border-[#5D8B57]/[0.35]
        p-[12px]
        flex flex-col
        pointer-events-auto
        z-[110]
        shadow-[0_12px_48px_rgba(0,0,0,0.5),0_0_30px_rgba(93,139,87,0.15),inset_0_1px_0_rgba(255,255,255,0.15)]
        before:absolute before:inset-0 before:rounded-[20px] before:bg-gradient-to-b before:from-white/[0.08] before:to-transparent before:pointer-events-none
      "
    >
      {/* TOP ICONS */}
      <div className="flex flex-col items-center gap-[16px]">
        {pages.map((page) => {
          const Icon = page.icon;
          const isActive = !isSettingsOpen && activeTab === page.id;
          const isDisabled = page.id === "mods";
          return (
            <button
              key={page.id}
              onClick={() => {
                if (isDisabled) return;
                console.log("Navbar click:", page.id);
                onTabChange(page.id);
              }}
              disabled={isDisabled}
              style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
              className={`relative transition-all duration-300 pointer-events-auto p-2 rounded-xl ${
                isDisabled
                  ? "cursor-not-allowed text-neutral-600 opacity-30"
                  : isActive 
                    ? "cursor-pointer text-[#7AB872] drop-shadow-[0_0_8px_rgba(122,184,114,0.8)] bg-[#5D8B57]/[0.2] shadow-[0_0_20px_rgba(93,139,87,0.3),inset_0_1px_0_rgba(255,255,255,0.2)]"
                    : "cursor-pointer text-white/70 hover:text-[#7AB872] hover:bg-[#5D8B57]/[0.1] hover:shadow-[0_0_15px_rgba(93,139,87,0.15)]"
              }`}
              title={isDisabled ? "Coming soon" : page.name}
            >
              <Icon size={18} />
            </button>
          );
        })} 
        {/* Divider */}
        <div
          className="w-[32px] h-[1px] bg-gradient-to-r from-transparent via-[#5D8B57]/40 to-transparent my-2"
          style={{ marginLeft: 'auto', marginRight: 'auto' }}
        />
        {/* Telegram icon */}
        <button
          type="button"
          onClick={openTelegram}
          style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
          className="transition-all w-[18px] h-[18px] cursor-pointer pointer-events-auto opacity-60 hover:opacity-90"
          title="Telegram"
        >
          <img src={telegramIcon} alt="Telegram" />
        </button>
        {/* Discord icon */}
        <button
          type="button"
          onClick={openDiscord}
          style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
          className="transition-all w-[18px] h-[18px] cursor-pointer pointer-events-auto opacity-60 hover:opacity-90"
          title="Discord"
        >
          <img src={discordIcon} alt="Discord" />
        </button>
        {/* Website icon */}
        <button
          type="button"
          onClick={openWebsite}
          style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
          className="transition-all w-[18px] h-[18px] cursor-pointer pointer-events-auto opacity-60 hover:opacity-90"
          title="Website"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-[18px] h-[18px] text-white">
            <circle cx="12" cy="12" r="10"/>
            <line x1="2" y1="12" x2="22" y2="12"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
          </svg>
        </button>
        <div
          className="w-[32px] h-[1px] bg-gradient-to-r from-transparent via-[#22d3ee]/30 to-transparent my-2"
          style={{ marginLeft: 'auto', marginRight: 'auto' }}
        />
        {/* Diagnostics icon */}
        <button
          type="button"
          onClick={onDiagnosticsClick}
          disabled={!onDiagnosticsClick}
          style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
          className={`transition-all cursor-not-allowed pointer-events-auto text-neutral-700 ${
            onDiagnosticsClick ? "opacity-30" : "opacity-60 hover:opacity-90"
          }`}
          title="Диагностика"
        >
          <Activity size={18} />
        </button>
        <button
          type="button"
          onClick={OpenFolder}
          style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
          className="transition-all cursor-pointer pointer-events-auto text-white opacity-60 hover:opacity-90"
          title="Папка игры"
        >
          <FolderOpen size={18} />
        </button>
        <button
          type="button"
          onClick={onSettingsClick}
          style={{ WebkitAppRegion: "no-drag" } as React.CSSProperties}
          className={`relative transition-all duration-300 cursor-pointer pointer-events-auto p-2 rounded-xl ${
            isSettingsOpen 
              ? "text-[#7AB872] drop-shadow-[0_0_8px_rgba(122,184,114,0.8)] bg-[#5D8B57]/[0.2] shadow-[0_0_20px_rgba(93,139,87,0.3),inset_0_1px_0_rgba(255,255,255,0.2)]"
              : "text-white/70 hover:text-[#7AB872] hover:bg-[#5D8B57]/[0.1] hover:shadow-[0_0_15px_rgba(93,139,87,0.15)]"
          }`}
          title="Настройки"
        >
          <Bolt size={18} />
        </button>
      </div>
    </div>
  );
}

export default Navbar;
