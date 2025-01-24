import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import Landing from '@/pages/Landing'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<Landing/>}/>
			</Routes>
		</BrowserRouter>
  </StrictMode>,
)
