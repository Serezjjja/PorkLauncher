import ParticleEffects from "./ParticleEffects";

function BackgroundImage() {
  return (
    <div className="absolute inset-0 z-0 scale-105 overflow-hidden">
      {/* Video Background */}
      <video
        autoPlay
        loop
        muted
        playsInline
        className="absolute inset-0 w-full h-full object-cover"
        style={{ objectFit: "cover" }}
      >
        <source src="backvid.mp4" type="video/mp4" />
      </video>

      {/* Magical atmosphere layers - Green Theme */}
      <div className="absolute inset-0 bg-gradient-to-br from-[#364E3A]/50 via-[#2A3D2E]/40 to-[#1C2A1F]/60" />

      {/* Green magical glow from top */}
      <div className="absolute inset-0 bg-gradient-to-b from-[#5D8B57]/15 via-transparent to-transparent" />

      {/* Green magical glow from corners */}
      <div className="absolute top-0 right-0 w-[600px] h-[600px] bg-[#7AB872]/12 rounded-full blur-[150px]" />
      <div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-[#5D8B57]/15 rounded-full blur-[120px]" />

      {/* Vignette effect */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_0%,rgba(13,18,30,0.4)_100%)]" />

      {/* Bottom fade for UI readability */}
      <div className="absolute inset-0 bg-gradient-to-t from-[#1C2A1F]/85 via-transparent to-transparent" />

      {/* Particle effects */}
      <ParticleEffects />
    </div>
  );
}

export default BackgroundImage;
