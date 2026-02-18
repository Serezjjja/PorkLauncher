import bgServersImage from "../assets/images/bg-servers.png";

import ParticleEffects from "./ParticleEffects";

function BackgroundServers() {
  return (
    <div
      className="absolute inset-0 bg-cover bg-center z-0 scale-105"
      style={{
        backgroundImage: `url(${bgServersImage})`,
      }}
    >
      {/* Magical atmosphere layers - Green Theme */}
      <div className="absolute inset-0 bg-gradient-to-br from-[#364E3A]/50 via-[#2A3D2E]/40 to-[#1C2A1F]/60" />
      
      {/* Green magical glow from top */}
      <div className="absolute inset-0 bg-gradient-to-b from-[#5D8B57]/15 via-transparent to-transparent" />
      
      {/* Green magical glow from corners */}
      <div className="absolute top-0 right-0 w-[600px] h-[600px] bg-[#7AB872]/12 rounded-full blur-[150px]" />
      <div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-[#5D8B57]/15 rounded-full blur-[120px]" />
      
      {/* Vignette effect */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_0%,rgba(28,42,31,0.4)_100%)]" />
      
      {/* Bottom fade for UI readability */}
      <div className="absolute inset-0 bg-gradient-to-t from-[#1C2A1F]/85 via-transparent to-transparent" />
      
      {/* Particle effects */}
      <ParticleEffects />
    </div>
  );
}

export default BackgroundServers;
