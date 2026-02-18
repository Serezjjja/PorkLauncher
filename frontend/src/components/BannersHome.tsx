import { useState } from "react";
import { AnimatePresence } from "framer-motion";
import { useTranslation } from "../i18n";
import { useLauncher, ServerWithFullUrls } from "../hooks/useLauncher";
import Banner from "./Banner";
import ServerModal from "./ServerModal";

interface BannersHomeProps {
  servers: ServerWithFullUrls[];
  isLoading: boolean;
  onPlay?: (serverIP: string) => void;
}

function BannersHome({ servers, isLoading, onPlay }: BannersHomeProps) {
  const { t } = useTranslation();
  const [selectedServer, setSelectedServer] = useState<ServerWithFullUrls | null>(null);

  // Show up to 5 banners
  const displayServers = servers.slice(0, 5);

  return (
    <div className="flex flex-col gap-[10px] max-h-[340px] overflow-y-auto scrollbar-hide">
      {/* Loading state - show placeholder banners (large size) */}
      {isLoading && (
        <>
          <div className="w-[448px] h-[200px] rounded-[20px] border border-[#5D8B57]/[0.25] bg-[#364E3A]/60 backdrop-blur-[24px] animate-pulse flex-shrink-0 shadow-[0_12px_48px_rgba(0,0,0,0.4),0_0_30px_rgba(93,139,87,0.12),inset_0_1px_0_rgba(255,255,255,0.1)]" />
        </>
      )}

      {/* Server banners - up to 5 (large variant like Servers page) */}
      {displayServers.map((server) => (
        <div key={server.id} className="flex-shrink-0">
          <Banner
            variant="large"
            backgroundImage={server.banner_url}
            iconImage={server.logo_url}
            title={server.name}
            description={server.description}
            onClick={() => setSelectedServer(server)}
          />
        </div>
      ))}

      {/* Fallback if no servers */}
      {!isLoading && displayServers.length === 0 && (
        <div className="w-[448px] h-[200px] rounded-[20px] border border-[#5D8B57]/[0.25] bg-[#364E3A]/60 backdrop-blur-[24px] flex items-center justify-center flex-shrink-0 shadow-[0_12px_48px_rgba(0,0,0,0.4),0_0_30px_rgba(93,139,87,0.12),inset_0_1px_0_rgba(255,255,255,0.1)]">
          <span className="text-[14px] text-white/50 font-[Mazzard] drop-shadow-[0_0_8px_rgba(255,255,255,0.2)]">
            {t.banners?.noServers || "No servers available"}
          </span>
        </div>
      )}

      {/* Server Detail Modal */}
      <AnimatePresence>
        {selectedServer && (
          <ServerModal
            server={selectedServer}
            isOpen={!!selectedServer}
            onClose={() => setSelectedServer(null)}
            onPlay={onPlay}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

export default BannersHome;
