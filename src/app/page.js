import Image from "next/image";
import Navbar from "./Navbar/page";
import Heropage from "./Hero/page"
import AboutPage from "./About/page";
import ContactPage from "./Contact/page";
export default function Home() {
  return(
<>
<Navbar>
  <AboutPage/>
  <ContactPage/>
</Navbar>
<Heropage></Heropage>

</>


  )
  
  ;
        
       
}
