import React, { useState } from "react";
import { motion } from "framer-motion";
import { X, Copy, Play, Check } from "lucide-react";
import { service } from "../../wailsjs/go/models";
import { useTranslation } from "../i18n";
import { useNavigateHome } from "../context/NavigationContext";

// Use the generated type
type ServerWithFullUrls = service.ServerWithUrls;

interface ServerModalProps {
  server: ServerWithFullUrls | null;
  isOpen: boolean;
  onClose: () => void;
  onPlay?: (serverIP: string) => void;
}

export const ServerModal: React.FC<ServerModalProps> = ({
  server,
  isOpen,
  onClose,
  onPlay,
}) => {
  const { t } = useTranslation();
  const navigateHome = useNavigateHome();
  const [copied, setCopied] = useState(false);

  if (!isOpen || !server) return null;

  const handleCopyIp = async () => {
    try {
      await navigator.clipboard.writeText(server.ip);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy IP:", err);
    }
  };

  const handlePlay = () => {
    if (onPlay) {
      onPlay(server.ip);
    }
    onClose();
    // Navigate to home page so user can see launch progress
    navigateHome();
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        transition={{ type: "spring", damping: 20, stiffness: 300 }}
        onClick={(e) => e.stopPropagation()}
        className="relative w-[480px] bg-[#364E3A]/80 backdrop-blur-[24px] rounded-[20px] border border-[#5D8B57]/[0.35] overflow-hidden shadow-[0_24px_64px_rgba(0,0,0,0.5),0_0_40px_rgba(93,139,87,0.15)]"
      >
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 z-10 p-2 text-white/50 hover:text-white transition-all duration-300 cursor-pointer bg-[#364E3A]/60 rounded-[10px] border border-[#5D8B57]/[0.3] hover:bg-[#5D8B57]/20 hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_15px_rgba(93,139,87,0.3)]"
        >
          <X size={18} strokeWidth={1.6} />
        </button>

        {/* Banner Image */}
        <div className="relative w-full h-[180px] overflow-hidden">
          {server.banner_url ? (
            <img
              src={server.banner_url}
              alt={server.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full bg-[#090909]/55 flex items-center justify-center">
              <span className="text-white/30 font-[Unbounded]">No Banner</span>
            </div>
          )}
          {/* Gradient overlay for smooth transition to content */}
          <div className="absolute bottom-0 left-0 right-0 h-[60px] bg-gradient-to-t from-[#364E3A]/80 to-transparent" />
        </div>

        {/* Content */}
        <div className="px-6 pb-6 -mt-6 relative">
          {/* Logo and Name Row */}
          <div className="flex items-start gap-4 mb-4">
            {/* Logo */}
            {server.logo_url && (
              <div className="w-[80px] h-[80px] rounded-[14px] overflow-hidden border border-[#5D8B57]/[0.3] bg-[#364E3A]/60 flex-shrink-0 shadow-[0_0_20px_rgba(93,139,87,0.15)]">
                <img
                  src={server.logo_url}
                  alt={`${server.name} logo`}
                  className="w-full h-full object-cover"
                />
              </div>
            )}

            {/* Name and IP */}
            <div className="flex-1 pt-2">
              <h2 className="text-[20px] font-[Unbounded] font-semibold text-white/90 tracking-[-0.02em] drop-shadow-[0_0_10px_rgba(93,139,87,0.4)]">
                {server.name}
              </h2>
              <p className="text-[14px] font-[Mazzard] text-[#7AB872] mt-1 drop-shadow-[0_0_6px_rgba(122,184,114,0.5)]">
                {server.ip}
              </p>
            </div>
          </div>

          {/* Description */}
          <p className="text-[15px] font-[Mazzard] text-white/70 leading-[1.5] mb-6">
            {server.description}
          </p>

          {/* Action Buttons */}
          <div className="flex gap-3">
            {/* Copy IP Button */}
            <button
              onClick={handleCopyIp}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-[#364E3A]/60 border border-[#5D8B57]/[0.3] rounded-[14px] hover:bg-[#5D8B57]/15 hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_20px_rgba(93,139,87,0.2)] transition-all duration-300 cursor-pointer group"
            >
              {copied ? (
                <>
                  <Check size={18} className="text-[#7AB872]" strokeWidth={1.6} />
                  <span className="text-[15px] font-[Mazzard] font-semibold text-[#7AB872] tracking-[-0.02em] drop-shadow-[0_0_6px_rgba(122,184,114,0.5)]">
                    {t.modals?.server?.copied || "Copied!"}
                  </span>
                </>
              ) : (
                <>
                  <Copy size={18} className="text-white/70 group-hover:text-[#7AB872] transition-colors duration-300" strokeWidth={1.6} />
                  <span className="text-[15px] font-[Mazzard] font-semibold text-white/70 group-hover:text-white tracking-[-0.02em]">
                    {t.modals?.server?.copyIp || "Copy IP"}
                  </span>
                </>
              )}
            </button>

            {/* Play Button */}
            <button
              onClick={handlePlay}
              className="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-[#5D8B57]/25 border border-[#5D8B57]/40 rounded-[14px] hover:bg-[#5D8B57]/40 hover:border-[#5D8B57]/60 hover:shadow-[0_0_25px_rgba(93,139,87,0.35)] transition-all duration-300 cursor-pointer group"
            >
              <Play size={18} className="text-[#7AB872] fill-[#7AB872] drop-shadow-[0_0_6px_rgba(122,184,114,0.8)]" strokeWidth={1.6} />
              <span className="text-[15px] font-[Mazzard] font-semibold text-[#7AB872] tracking-[-0.02em] drop-shadow-[0_0_6px_rgba(122,184,114,0.5)]">
                {t.modals?.server?.play || "Play"}
              </span>
            </button>
          </div>
        </div>
      </motion.div>
    </motion.div>
  );
};

export default ServerModal;
