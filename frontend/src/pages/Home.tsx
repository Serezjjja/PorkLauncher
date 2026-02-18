import React, { useContext } from "react";
import { AnimatePresence, motion } from "framer-motion";
import BannersHome from "../components/BannersHome";
import { ControlSection } from "../components/ControlSection";
import { ProfileSection } from "../components/ProfileCard";
import { DiscordWidget } from "../components/DiscordWidget";
import { TelegramWidget } from "../components/TelegramWidget";
import { UpdateOverlay } from "../components/UpdateOverlay";
import { DeleteConfirmationModal } from "../components/DeleteConfirmationModal";
import { ErrorModal } from "../components/ErrorModal";
import { LoginModal } from "../components/LoginModal";
import { useLauncher } from "../hooks/useLauncher";
import { OpenFolder, DeleteGame } from "../../wailsjs/go/app/App";
import { SettingsOverlayContext } from "../context/SettingsOverlayContext";

function HomePage() {
  const showSettingsOverlay = useContext(SettingsOverlayContext);
  const {
    username,
    currentVersion,
    selectedBranch,
    availableVersions,
    isLoadingVersions,
    launcherVersion,
    isEditingUsername,
    setIsEditingUsername,
    progress,
    status,
    isDownloading,
    downloadDetails,
    updateAsset,
    isUpdatingLauncher,
    updateStats,
    showDeleteModal,
    setShowDeleteModal,
    setShowDiagnostics,
    error,
    setError,
    handlePlay,
    handleUpdateLauncher,
    setNick,
    setLocalGameVersion,
    handleBranchChange,
    servers,
    isLoadingServers,
    // Auth
    isAuthenticated,
    showLoginModal,
    setShowLoginModal,
    handleLogin,
    handleLogout,
    isAuthLoading,
    authError,
  } = useLauncher();

  return (
    <>
      {isUpdatingLauncher && (
        <UpdateOverlay
          progress={progress}
          downloaded={updateStats.d}
          total={updateStats.t}
        />
      )}

      <main className="relative z-10 h-full p-10 flex flex-col justify-between pt-[60px]">
        <div className="flex justify-between items-start">
          <div className="flex flex-col gap-4">
            <ProfileSection
              username={username}
              currentVersion={currentVersion}
              selectedBranch={selectedBranch}
              availableVersions={availableVersions}
              isLoadingVersions={isLoadingVersions}
              isEditing={isEditingUsername}
              onEditToggle={setIsEditingUsername}
              onUserChange={setNick}
              onVersionChange={setLocalGameVersion}
              onBranchChange={handleBranchChange}
              isAuthenticated={isAuthenticated}
              onLogout={handleLogout}
            />
          </div>
          <AnimatePresence mode="wait">
            {!showSettingsOverlay && (
              <motion.div
                key="banners-home"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.28, ease: [0.16, 1, 0.3, 1] }}
                className="flex flex-col gap-4"
              >
                <BannersHome servers={servers} isLoading={isLoadingServers} onPlay={() => handlePlay()} />
                <div className="flex gap-4">
                  <DiscordWidget />
                  <TelegramWidget />
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        <ControlSection
          onPlay={() => handlePlay()}
          isDownloading={isDownloading}
          progress={progress}
          status={status}
          speed={downloadDetails.speed}
          downloaded={downloadDetails.downloaded}
          total={downloadDetails.total}
          currentFile={downloadDetails.currentFile}
          actions={{
            openFolder: OpenFolder,
            showDiagnostics: () => setShowDiagnostics(true),
            showDelete: () => setShowDeleteModal(true),
          }}
          updateAvailable={!!updateAsset}
          onUpdate={handleUpdateLauncher}
        />
      </main>

      {showDeleteModal && (
        <DeleteConfirmationModal
          onConfirm={() => {
            DeleteGame("default");
            setShowDeleteModal(false);
          }}
          onCancel={() => setShowDeleteModal(false)}
        />
      )}

      {error && <ErrorModal error={error} onClose={() => setError(null)} />}

      {/* Login Modal */}
      <LoginModal
        isOpen={showLoginModal}
        onClose={() => setShowLoginModal(false)}
        onLogin={handleLogin}
        error={authError}
        isLoading={isAuthLoading}
      />
    </>
  );
}

export default HomePage;