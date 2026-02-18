import React, { useMemo, useRef } from "react";
import { SquarePen, LogOut, User } from "lucide-react";
import { useTranslation } from "../i18n";

export type ReleaseType = "pre-release" | "release";

interface VersionOption {
  value: string;
  label: string;
}

interface ProfileProps {
  username: string;
  currentVersion: string;
  selectedBranch: ReleaseType;
  availableVersions: VersionOption[];
  isLoadingVersions?: boolean;
  isEditing: boolean;
  onEditToggle: (val: boolean) => void;
  onUserChange: (val: string) => void;
  onVersionChange: (val: string) => void;
  onBranchChange: (branch: ReleaseType) => void;
  // Auth props
  isAuthenticated?: boolean;
  onLogout?: () => void;
}

export const ProfileSection: React.FC<ProfileProps> = ({
  username,
  currentVersion,
  selectedBranch,
  availableVersions,
  isLoadingVersions = false,
  isEditing,
  onEditToggle,
  onUserChange,
  onVersionChange,
  onBranchChange,
  isAuthenticated = false,
  onLogout,
}) => {
  const { t } = useTranslation();

  // Only Release version is available (pre-release removed)
  const OPTIONS: { value: ReleaseType; label: string }[] = [
    { 
      value: "release", 
      label: t.profile.releaseType.release || "Release" 
    },
  ];

  const rootRef = useRef<HTMLDivElement | null>(null);

  const baseText = "text-white/90 font-[MazzardM-Medium] text-[16px] drop-shadow-[0_0_8px_rgba(255,255,255,0.3)]";
  const glass =
    "bg-[#364E3A]/[0.75] backdrop-blur-[24px] border border-[#5D8B57]/[0.35] shadow-[0_12px_48px_rgba(0,0,0,0.5),0_0_30px_rgba(93,139,87,0.15),inset_0_1px_0_rgba(255,255,255,0.15)]";
  const hover = "hover:bg-[#5D8B57]/[0.12] hover:border-[#5D8B57]/[0.5] hover:shadow-[0_0_20px_rgba(93,139,87,0.25)] transition-all duration-300";

  // Get display label for current branch
  const currentBranchLabel = useMemo(() => {
    const option = OPTIONS.find(opt => opt.value === selectedBranch);
    return option?.label || "Release";
  }, [selectedBranch, OPTIONS]);

  return (
    <div className="ml-[48px]" ref={rootRef}>
      {/* Username with Auth Status */}
      <div
        className={`w-[280px] h-[52px] ${glass} rounded-[16px] p-4 flex items-center justify-between mb-3 relative overflow-hidden before:absolute before:inset-0 before:bg-gradient-to-r before:from-transparent before:via-white/[0.03] before:to-transparent before:pointer-events-none`}
      >
        {isEditing ? (
          <input
            autoFocus
            className={`${baseText} bg-transparent outline-none tracking-[-3%] w-full`}
            defaultValue={username}
            onBlur={(e: React.FocusEvent<HTMLInputElement>) => {
              onEditToggle(false);
              onUserChange(e.target.value);
            }}
            onKeyDown={(e: React.KeyboardEvent) =>
              e.key === "Enter" && (e.target as HTMLInputElement).blur()
            }
          />
        ) : (
          <>
            <div className="flex items-center gap-2">
              <User
                size={18}
                className={`${isAuthenticated ? "text-[#7AB872] drop-shadow-[0_0_8px_rgba(122,184,114,0.8)]" : "text-[#5D8B57] drop-shadow-[0_0_6px_rgba(93,139,87,0.6)]"}`}
              />
              <span className={baseText}>{username}</span>
              {isAuthenticated && (
                <span className="text-[10px] font-[Mazzard] text-[#7AB872] bg-[#5D8B57]/20 px-2 py-0.5 rounded-full shadow-[0_0_10px_rgba(93,139,87,0.3)] border border-[#5D8B57]/40">
                  ✓
                </span>
              )}
            </div>
            <div className="flex items-center gap-2">
              {!isAuthenticated && (
                <SquarePen
                  size={16}
                  className="text-white/60 hover:text-[#7AB872] cursor-pointer w-[16px] h-[16px] transition-all duration-300 hover:drop-shadow-[0_0_6px_rgba(122,184,114,0.8)]"
                  onClick={() => onEditToggle(true)}
                />
              )}
              {isAuthenticated && onLogout && (
                <button
                  onClick={onLogout}
                  className="p-1.5 text-white/50 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-all duration-300 hover:shadow-[0_0_10px_rgba(248,113,113,0.3)]"
                  title={t.auth?.logout || "Logout"}
                >
                  <LogOut size={14} />
                </button>
              )}
            </div>
          </>
        )}
      </div>

      {/* Bottom pill - Release indicator only */}
      <div
        className={`relative w-[140px] h-[52px] ${glass} rounded-[16px] flex overflow-hidden before:absolute before:inset-0 before:bg-gradient-to-r before:from-transparent before:via-white/[0.03] before:to-transparent before:pointer-events-none`}
      >
        {/* Release type indicator */}
        <div
          className={`
            relative w-full h-full px-[16px] 
            flex items-center justify-between 
            rounded-[16px]
          `}
        >
          <span className={`${baseText} truncate text-[#7AB872]`}>{currentBranchLabel}</span>
          <div className="w-2 h-2 rounded-full bg-[#7AB872] shadow-[0_0_8px_rgba(122,184,114,0.8)]" />
        </div>
      </div>


    </div>
  );
};