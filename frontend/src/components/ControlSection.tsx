import React from "react";
import { motion } from "framer-motion";
import {
  FolderOpen,
  Activity,
  Settings,
  Trash,
  ArrowUpCircle,
  RefreshCcw,
} from "lucide-react";
import { useTranslation } from "../i18n";

interface ControlSectionProps {
  onPlay: () => void;
  isDownloading: boolean;
  progress: number;
  status: string;
  speed: string; // Added
  downloaded: number; // Added
  total: number; // Added
  currentFile: string; // Added
  actions: {
    openFolder: () => void;
    showDiagnostics: () => void;
    showDelete: () => void;
  };
  updateAvailable: boolean;
  onUpdate: () => void;
}

export const ControlSection: React.FC<ControlSectionProps> = ({
  onPlay,
  isDownloading,
  progress,
  status,
  speed,
  downloaded,
  total,
  currentFile,
  actions,
  updateAvailable,
  onUpdate,
}) => {
  const { t } = useTranslation();

  // Your original formatting helper
  const formatBytes = (bytes: number) => {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  return (
    <div className="w-full flex items-end gap-[24px] ml-[48px]">
      <div className="w-[280px] flex flex-col gap-[12px]">
        {updateAvailable && (
          <button
            onClick={onUpdate}
            className="cursor-pointer hover:scale-102 w-[280px] h-[44px] bg-[#364E3A]/[0.75] backdrop-blur-[24px] border border-[#7AB872]/[0.4] rounded-[16px] px-[16px] flex items-center justify-between shadow-[0_8px_32px_rgba(0,0,0,0.4),0_0_20px_rgba(122,184,114,0.15),inset_0_1px_0_rgba(255,255,255,0.1)] hover:shadow-[0_8px_32px_rgba(0,0,0,0.4),0_0_30px_rgba(122,184,114,0.25),inset_0_1px_0_rgba(255,255,255,0.15)] transition-all duration-300 group"
          >
            <span className="text-[16px] text-white/90 font-[Mazzard] tracking-[-3%] drop-shadow-[0_0_8px_rgba(255,255,255,0.2)]">
              {t.control.updateAvailable}
            </span>
            <span className="text-[#7AB872] transition-transform flex items-center justify-center group-hover:rotate-180 duration-500">
              <RefreshCcw size={18} className="drop-shadow-[0_0_6px_rgba(122,184,114,0.8)]" />
            </span>
          </button>
        )}
        <motion.button
          style={{
            willChange: "transform",
            transform: "translateZ(0)",
            backfaceVisibility: "hidden",
          }}
          whileHover={isDownloading ? {} : { 
            scale: 1.03,
            boxShadow: "0 0 80px rgba(93,139,87,0.8), 0 0 120px rgba(122,184,114,0.5), inset 0 2px 0 rgba(255,255,255,0.5)"
          }}
          whileTap={isDownloading ? {} : { scale: 0.95 }}
          animate={!isDownloading ? {
            boxShadow: [
              "0 0 40px rgba(93,139,87,0.4), 0 0 80px rgba(122,184,114,0.25), inset 0 2px 0 rgba(255,255,255,0.3)",
              "0 0 60px rgba(93,139,87,0.6), 0 0 100px rgba(122,184,114,0.4), inset 0 2px 0 rgba(255,255,255,0.35)",
              "0 0 40px rgba(93,139,87,0.4), 0 0 80px rgba(122,184,114,0.25), inset 0 2px 0 rgba(255,255,255,0.3)"
            ]
          } : {}}
          transition={{
            type: "spring",
            stiffness: 400,
            damping: 25,
            boxShadow: {
              duration: 2,
              repeat: Infinity,
              ease: "easeInOut"
            }
          }}
          onClick={onPlay}
          disabled={isDownloading}
          className={`relative w-[280px] h-[100px] font-[Unbounded] font-[700] text-[28px] text-white overflow-hidden rounded-[50px] border border-[#5D8B57]/[0.5] shadow-[0_0_40px_rgba(93,139,87,0.4),0_0_80px_rgba(122,184,114,0.25),inset_0_2px_0_rgba(255,255,255,0.3)] disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 group ${
            isDownloading ? "cursor-not-allowed" : "cursor-pointer"
          }`}
        >
          {/* Animated gradient background - GREEN theme */}
          <motion.div 
            className="absolute inset-0 bg-gradient-to-b from-[#5D8B57] via-[#4A7A44] to-[#364E3A] opacity-95"
            animate={!isDownloading ? {
              background: [
                "linear-gradient(to bottom, #5D8B57, #4A7A44, #364E3A)",
                "linear-gradient(to bottom, #6B9B65, #5A8A54, #465E4A)",
                "linear-gradient(to bottom, #5D8B57, #4A7A44, #364E3A)"
              ]
            } : {}}
            transition={{
              duration: 3,
              repeat: Infinity,
              ease: "easeInOut"
            }}
          />
          
          {/* Animated side glow - GREEN */}
          <motion.div 
            className="absolute inset-0"
            animate={!isDownloading ? {
              background: [
                "linear-gradient(to right, rgba(122,184,114,0.4), transparent, rgba(93,139,87,0.4))",
                "linear-gradient(to right, rgba(93,139,87,0.6), transparent, rgba(122,184,114,0.6))",
                "linear-gradient(to right, rgba(122,184,114,0.4), transparent, rgba(93,139,87,0.4))"
              ]
            } : {}}
            transition={{
              duration: 2.5,
              repeat: Infinity,
              ease: "easeInOut"
            }}
          />
          
          {/* Glossy highlight with shimmer */}
          <div className="absolute inset-x-0 top-0 h-[45%] bg-gradient-to-b from-white/40 to-transparent rounded-t-[50px]" />
          
          {/* Moving shimmer effect */}
          <motion.div 
            className="absolute inset-0 opacity-40"
            style={{
              background: "linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.3) 50%, transparent 100%)",
              backgroundSize: "200% 100%"
            }}
            animate={!isDownloading ? {
              backgroundPosition: ["200% 0%", "-200% 0%"]
            } : {}}
            transition={{
              duration: 3,
              repeat: Infinity,
              ease: "linear"
            }}
          />
          
          {/* Pulsing inner glow */}
          <motion.div
            className="absolute inset-2 rounded-[46px] border border-[#7AB872]/[0.3]"
            animate={!isDownloading ? {
              opacity: [0.3, 0.6, 0.3],
              scale: [1, 1.02, 1]
            } : {}}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: "easeInOut"
            }}
          />
          
          {/* Text with enhanced glow animation */}
          <motion.span 
            className="relative z-10 tracking-wider"
            animate={!isDownloading ? {
              textShadow: [
                "0 2px 4px rgba(0,0,0,0.3), 0 0 20px rgba(122,184,114,0.5)",
                "0 2px 4px rgba(0,0,0,0.3), 0 0 40px rgba(122,184,114,0.8), 0 0 60px rgba(122,184,114,0.4)",
                "0 2px 4px rgba(0,0,0,0.3), 0 0 20px rgba(122,184,114,0.5)"
              ]
            } : {
              textShadow: "0 2px 4px rgba(0,0,0,0.3)"
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: "easeInOut"
            }}
          >
            {isDownloading ? t.common.install : t.common.play}
          </motion.span>
        </motion.button>
      </div>

      <div className="flex-1 flex flex-col gap-[8px] pb-2">
        <div className="flex justify-between items-end">
          <div className="tracking-[-3%] flex items-baseline gap-[20px]">
            <span className="text-[36px] text-white font-[Unbounded] font-[600] drop-shadow-[0_0_15px_rgba(34,211,238,0.5)]">
              {Math.round(progress)}%
            </span>
            <span className="text-[16px] text-white/50 font-[Mazzard] drop-shadow-[0_0_8px_rgba(255,255,255,0.2)]">
              {status}
            </span>
          </div>

          {/* Re-added speed and total size labels */}
          <div className="text-[14px] text-white/40 font-[MazzardM-Medium] text-right break-words min-w-0 flex-1 mr-[48px] drop-shadow-[0_0_6px_rgba(255,255,255,0.15)]">
            {speed && total > 0
              ? `${speed} • ${formatBytes(downloaded)} / ${formatBytes(total)}`
              : currentFile || t.common.ready}
          </div>
        </div>
        
        {/* Energy Beam Progress Bar */}
        <div className="h-[10px] w-[852px] bg-[#364E3A]/[0.6] rounded-full overflow-hidden border border-[#5D8B57]/[0.25] shadow-[inset_0_2px_4px_rgba(0,0,0,0.3),0_0_20px_rgba(93,139,87,0.15)] relative">
          {/* Background glow */}
          <div className="absolute inset-0 bg-gradient-to-r from-[#5D8B57]/8 via-[#7AB872]/12 to-[#5D8B57]/8" />
          
          <motion.div
            style={{
              willChange: "width",
              transform: "translateZ(0)",
              backfaceVisibility: "hidden",
            }}
            animate={{ width: `${progress}%` }}
            transition={{
              type: "tween",
              ease: "linear",
              duration: 0.1,
            }}
            className="h-full relative energy-beam"
          >
            {/* Leading edge glow */}
            <div className="absolute right-0 top-1/2 -translate-y-1/2 w-[20px] h-[20px] bg-[#7AB872] rounded-full blur-[8px] opacity-80" />
            <div className="absolute right-0 top-1/2 -translate-y-1/2 w-[8px] h-[8px] bg-white rounded-full shadow-[0_0_15px_#fff,0_0_30px_#7AB872]" />
          </motion.div>
        </div>
      </div>
    </div>
  );
};

const NavBtn = ({ icon, onClick }: { icon: any; onClick?: () => void }) => (
  <button
    onClick={onClick}
    className="w-[66px] h-[44px] cursor-pointer flex items-center justify-center bg-[#364E3A]/[0.75] backdrop-blur-[24px] border border-[#5D8B57]/[0.3] rounded-[16px] hover:bg-[#5D8B57]/[0.12] hover:border-[#5D8B57]/[0.55] hover:shadow-[0_0_20px_rgba(93,139,87,0.25)] transition-all duration-300 text-white/60 hover:text-[#7AB872] hover:drop-shadow-[0_0_8px_rgba(122,184,114,0.8)]"
  >
    {icon}
  </button>
);
