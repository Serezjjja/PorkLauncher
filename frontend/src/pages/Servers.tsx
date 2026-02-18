import { useState } from "react";
import { AnimatePresence } from "framer-motion";
import Banner from "../components/Banner";
import { useTranslation } from "../i18n";
import { useLauncher, ServerWithFullUrls } from "../hooks/useLauncher";
import ServerModal from "../components/ServerModal";

function ServersPage() {
  const { t } = useTranslation();
  const { servers, isLoadingServers, handlePlay } = useLauncher();
  const [selectedServer, setSelectedServer] = useState<ServerWithFullUrls | null>(null);

  // Show all servers
  const displayServers = servers;

  return (
    <div className="relative h-full w-full">
      {/* Title */}
      <div
        className="
          absolute
          left-[88px]
          top-[58px]
          text-white/90
          text-[22px]
          font-[600]
          tracking-[0.04em]
          uppercase
          font-[Unbounded]
          drop-shadow-[0_0_15px_rgba(93,139,87,0.5)]
        "
      >
        {t.pages.servers}
      </div>

      {/* Loading State */}
      {isLoadingServers && (
        <div className="absolute left-[88px] top-[100px] flex flex-wrap gap-x-[22px] gap-y-[22px]">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="w-[448px] h-[200px] rounded-[20px] bg-[#364E3A]/60 backdrop-blur-[24px] animate-pulse border border-[#5D8B57]/[0.25] shadow-[0_12px_48px_rgba(0,0,0,0.4),0_0_30px_rgba(93,139,87,0.12),inset_0_1px_0_rgba(255,255,255,0.1)]" />
          ))}
        </div>
      )}

      {/* Servers Grid */}
      <div className="absolute left-[88px] top-[100px] flex flex-wrap gap-x-[22px] gap-y-[22px]">
        {/* Large Server Banners */}
        {displayServers.map((server) => (
          <Banner
            key={server.id}
            variant="large"
            backgroundImage={server.banner_url}
            iconImage={server.logo_url}
            title={server.name}
            description={server.description}
            onClick={() => setSelectedServer(server)}
          />
        ))}

      </div>

      {/* Server Detail Modal */}
      <AnimatePresence>
        {selectedServer && (
          <ServerModal
            server={selectedServer}
            isOpen={!!selectedServer}
            onClose={() => setSelectedServer(null)}
            onPlay={handlePlay}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

export default ServersPage;
