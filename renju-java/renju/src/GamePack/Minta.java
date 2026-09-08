package GamePack;

import java.util.ArrayList;
import java.util.List;
import java.util.Vector;

public class Minta {

    public Vector<String> whiteFive;
    public Vector<String> blackFive;
    public Vector<String> overLine;
    public Vector<String> inlineDouble;
    public Vector<String> fours;
    public Vector<String> threes;
    public ArrayList<Vector<String>> mintalista;

    //Kontruktor: Listaba fuzi a tiltott lepesek es nyertes lepesekre mintait.
    //Azutan pedig a tiltott lepesek fajtaibol csinal meg egy listat
    public void Minta() {
        whiteFive = new Vector<>(List.of("OOOOO"));
        blackFive = new Vector<>(List.of("XXXXX"));
        overLine = new Vector<>(List.of("XXXXXX"));
        inlineDouble = new Vector<>(List.of("-XX-XX-XX-", "-XX-XX-XXO", "-XX-XX-XX#",
                                            "OXX-XX-XX-", "OXX-XX-XXO", "OXX-XX-XX#",
                                            "#XX-XX-XX-", "#XX-XX-XXO", "#XX-XX-XX#",
                                            "#X-XXX-X-", "#X-XXX-XO", "#X-XXX-XO",
                                            "#X-XXX-X-", "#X-XXX-XO", "#X-XXX-X#"));
        fours = new Vector<>(List.of("XXXX-", "-XXXX", "XXX-X", "XX-XX", "X-XXX"));
        threes = new Vector<>(List.of("--XXX-", "-XXX--", "-XX-X-", "-X-XX-"));

        mintalista = new ArrayList<>(List.of(whiteFive, blackFive, overLine, inlineDouble, fours, threes));
    }

    //Megkapja egy minta lista meretet
    public int mintaListaMeret(int type) {
        int result = mintalista.get(type).size();
        return result;
    }

    //Megkapja egy minta listanak egy specifikus stringjet
    public String getMinta(int type, int idx) {
        String minta = new String();
        minta = mintalista.get(type).get(idx);
        return minta;
    }
}
