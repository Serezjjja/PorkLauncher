import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { X, Loader2, User, Lock, AlertCircle } from "lucide-react";
import { useTranslation } from "../i18n";

interface LoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLogin: (email: string, password: string) => Promise<void>;
  error: string | null;
  isLoading: boolean;
}

export const LoginModal: React.FC<LoginModalProps> = ({
  isOpen,
  onClose,
  onLogin,
  error,
  isLoading,
}) => {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim() || !password.trim() || isLoading) return;
    await onLogin(email, password);
  };

  // Reset form when modal closes
  React.useEffect(() => {
    if (!isOpen) {
      setEmail("");
      setPassword("");
    }
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-[100] p-4"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.9, opacity: 0 }}
            transition={{ type: "spring", damping: 20, stiffness: 300 }}
            onClick={(e) => e.stopPropagation()}
            className="relative w-[420px] bg-[#364E3A]/80 backdrop-blur-[24px] rounded-[20px] border border-[#5D8B57]/[0.35] overflow-hidden shadow-[0_24px_64px_rgba(0,0,0,0.5),0_0_40px_rgba(93,139,87,0.15)]"
          >
            {/* Close Button */}
            <button
              onClick={onClose}
              className="absolute top-4 right-4 z-10 p-2 text-white/50 hover:text-white transition-all duration-300 cursor-pointer bg-[#364E3A]/60 rounded-[10px] border border-[#5D8B57]/[0.3] hover:bg-[#5D8B57]/20 hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_15px_rgba(93,139,87,0.3)]"
            >
              <X size={18} strokeWidth={1.6} />
            </button>

            {/* Header */}
            <div className="px-8 pt-8 pb-4">
              <h2 className="text-[24px] font-[Unbounded] font-semibold text-white/90 tracking-[-0.02em] drop-shadow-[0_0_10px_rgba(93,139,87,0.4)]">
                {t.auth?.loginTitle || "Login"}
              </h2>
              <p className="text-[14px] font-[Mazzard] text-white/60 mt-2">
                {t.auth?.loginSubtitle || "Sign in with your account to play"}
              </p>
            </div>

            {/* Form */}
            <form onSubmit={handleSubmit} className="px-8 pb-8 space-y-4">
              {/* Error Message */}
              {error && (
                <motion.div
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="flex items-center gap-2 p-3 bg-red-500/10 border border-red-500/20 rounded-[12px] text-red-400 text-[13px] font-[Mazzard]"
                >
                  <AlertCircle size={16} />
                  <span>{error}</span>
                </motion.div>
              )}

              {/* Email Field */}
              <div className="space-y-2">
                <label className="text-[13px] font-[Mazzard] font-medium text-[#7AB872] ml-1 drop-shadow-[0_0_4px_rgba(122,184,114,0.3)]">
                  {t.auth?.emailLabel || "Email"}
                </label>
                <div className="relative">
                  <User
                    size={18}
                    className="absolute left-4 top-1/2 -translate-y-1/2 text-[#5D8B57]"
                  />
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder={t.auth?.emailPlaceholder || "your@email.com"}
                    className="w-full h-[48px] pl-12 pr-4 bg-[#364E3A]/50 border border-[#5D8B57]/[0.25] rounded-[14px] text-white/90 font-[Mazzard] text-[15px] placeholder:text-white/30 focus:outline-none focus:border-[#5D8B57]/[0.6] focus:bg-[#5D8B57]/[0.08] focus:shadow-[0_0_20px_rgba(93,139,87,0.15)] transition-all duration-300"
                    disabled={isLoading}
                  />
                </div>
              </div>

              {/* Password Field */}
              <div className="space-y-2">
                <label className="text-[13px] font-[Mazzard] font-medium text-[#7AB872] ml-1 drop-shadow-[0_0_4px_rgba(122,184,114,0.3)]">
                  {t.auth?.passwordLabel || "Password"}
                </label>
                <div className="relative">
                  <Lock
                    size={18}
                    className="absolute left-4 top-1/2 -translate-y-1/2 text-[#5D8B57]"
                  />
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder={t.auth?.passwordPlaceholder || "••••••••"}
                    className="w-full h-[48px] pl-12 pr-4 bg-[#364E3A]/50 border border-[#5D8B57]/[0.25] rounded-[14px] text-white/90 font-[Mazzard] text-[15px] placeholder:text-white/30 focus:outline-none focus:border-[#5D8B57]/[0.6] focus:bg-[#5D8B57]/[0.08] focus:shadow-[0_0_20px_rgba(93,139,87,0.15)] transition-all duration-300"
                    disabled={isLoading}
                  />
                </div>
              </div>

              {/* Submit Button */}
              <motion.button
                type="submit"
                disabled={isLoading || !email.trim() || !password.trim()}
                whileTap={{ scale: 0.97 }}
                className="relative w-full h-[52px] mt-4 overflow-hidden rounded-[14px] border border-[#5D8B57]/[0.5] shadow-[0_0_30px_rgba(93,139,87,0.25)] hover:shadow-[0_0_40px_rgba(93,139,87,0.4)] transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 group"
              >
                {/* Gradient background */}
                <div className="absolute inset-0 bg-gradient-to-r from-[#5D8B57] via-[#4A7A44] to-[#5D8B57] opacity-90" />
                <div className="absolute inset-0 bg-gradient-to-b from-[#7AB872]/20 to-transparent" />
                
                {/* Glossy highlight */}
                <div className="absolute inset-x-0 top-0 h-[40%] bg-gradient-to-b from-white/30 to-transparent" />
                
                {/* Shimmer effect */}
                <div className="absolute inset-0 shimmer opacity-20" />
                
                {/* Content */}
                {isLoading ? (
                  <>
                    <Loader2 size={18} className="animate-spin relative z-10 text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.3)]" />
                    <span className="relative z-10 text-white font-[Unbounded] font-semibold text-[15px] tracking-[-0.01em] drop-shadow-[0_2px_4px_rgba(0,0,0,0.3)]">{t.auth?.loggingIn || "Signing in..."}</span>
                  </>
                ) : (
                  <span className="relative z-10 text-white font-[Unbounded] font-semibold text-[15px] tracking-[-0.01em] drop-shadow-[0_2px_4px_rgba(0,0,0,0.3)]">{t.auth?.loginButton || "Sign In"}</span>
                )}
              </motion.button>

              {/* Register Link */}
              <div className="text-center pt-2">
                <span className="text-[13px] font-[Mazzard] text-white/50">
                  {t.auth?.noAccount || "Don't have an account?"}{" "}
                </span>
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    // Open registration page in browser
                    window.open("https://porkland.net/user/register", "_blank");
                  }}
                  className="text-[13px] font-[Mazzard] text-[#7AB872] hover:text-[#5D8B57] transition-colors duration-300 drop-shadow-[0_0_6px_rgba(122,184,114,0.4)]"
                >
                  {t.auth?.registerLink || "Register"}
                </a>
              </div>
            </form>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default LoginModal;
