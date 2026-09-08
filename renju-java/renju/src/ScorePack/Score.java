package ScorePack;

import java.io.Serializable;

public class Score implements Serializable {

   // private int sorszam; i dnt know man

    private String lepesszam;

    private String feherlepes;

    private String feketelepes;

    private String nyertes;

    // Osszlepesszam lekerdezes es beallitasa
    public String getLepesszam(){
        return lepesszam;
    }
    public void setLepesszam(String lepesszam){
        this.lepesszam = lepesszam;
    }

    // Fehere lepesszama lekerdezese es beallitasa
    public String getFeherlepes(){
        return feherlepes;
    }
    public void setFeherlepes(String feherlepes){
        this.feherlepes = feherlepes;
    }

    // FEKETE lepesszam lekeredeezese es beallitasa
    public String getFeketelepes(){
        return feketelepes;
    }
    public void setFeketelepes(String feketelepes){
        this.feketelepes = feketelepes;
    }

    // Nyeretes lekerdezese es beallitasa
    public String getNyertes(){
        return nyertes;
    }
    public void setNyertes( String nyertes){
        this.nyertes = nyertes;
    }

}
