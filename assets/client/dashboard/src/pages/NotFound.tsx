import "@/assets/notfound.css";

function NotFound() {
  return (
    <>
      <div className="error-container">
        <div className="error-content">
          <div className="moon-illustration">
            <div className="crater crater-1"></div>
            <div className="crater crater-2"></div>
            <div className="crater crater-3"></div>
          </div>

          <div className="error-code">404</div>
          <h1 className="error-title">Page not found</h1>
          <p className="error-message">
            It seems like you've ventured into unknown space. The page you're
            looking for doesn't exist or has been moved to another location.
          </p>
        </div>
      </div>
    </>
  );
}

export default NotFound;
