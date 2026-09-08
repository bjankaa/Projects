package GamePack;


import java.io.*;
import java.util.Objects;

public class Game {

    private String color;
    private int stone_x;
    private int stone_y;
    private String[][] tabla;
    private Pair<Integer, String>[] axes;
    private Minta lista ;

    // "Mozgatashoz" szukseges matrix
    private int[][] increment = {{-1, 0},       //felfele
                                {-1, 1},      //jobbra fel
                                {0, 1},        //jobbra oldalra
                                {1, 1}};        //jobbra le

    // Pair tipusu tomb initje
    public Pair<Integer, String>[] axesinit() {
        this.axes = new Pair[4];
        String ures = "";
        for (int i = 0; i < 4; i++) {
            axes[i] = new Pair<>(i, ures);
        }
        return axes;
    }

    // gettes és setter rengeteg.Fokepp tesztekhez kellenek a getterek
    public void setStone_x(int stone_x) {
        this.stone_x = stone_x;
    }
    public int getStone_x() {
        return this.stone_x;
    }

    public void setStone_y(int stone_y) {
        this.stone_y = stone_y;
    }
    public int getStone_y() {
        return this.stone_y;
    }

    public void setTabla(String[][] tomb) {
        this.tabla = tomb;
    }
    public String[][] getTabla() {
        return this.tabla;
    }

    public void setColor(String color) {
        this.color = color;
    }
    public String getColor() {
        return this.color;
    }

    // Jelenlegi kovunk stringjeit kulonbozo axisokra nezve lementjuk
    public void setStone(int row,  int col){
        for(int i =0; i< 4; i++){
            int currentpos = getPosition(row, col, i);
            String sqnc = getSequence(row, col, i, currentpos);
            axes[i] = new Pair<>(currentpos,sqnc);
        }
    }


    // Currentaxis azt valasztja ki merre fogunk lepkedni, mikor megnezuk milyen tavolsagra vagyunk a szeletol
    //ebben a fuggvenyben megkapjuk iranytol fuggoen merre vagyunk. Ellenorzeseknel kell ezt tudnunk.
    public int getPosition(int row, int col, int currentAxis) {

        int positionidx = 0;

        while (tabla[row][col] != "#") {
            positionidx++;
            row -= increment[currentAxis][0];
            col -= increment[currentAxis][1];
        }
        return positionidx;
    }


    //Az adott tengelyben ahol a kovunk van abbol egy stringet alkot. Ebben a stringben fogjuk kesobb meg keresni a mintak az ellenorzesnel
    String getSequence(int row, int col, int currentAxis, int currentposition) {

        String axis_Sqnc = "";

        row -= ((currentposition - 1) * increment[currentAxis][0]);
        col -= ((currentposition - 1) * increment[currentAxis][1]);

        axis_Sqnc += "#";

        while (!tabla[row][col].equals("#")) {
            axis_Sqnc += tabla[row][col];
            row += increment[currentAxis][0];
            col += increment[currentAxis][1];
        }

        axis_Sqnc += "#";
        return axis_Sqnc;

    }

    // Adott axison megkapjuk kovunk stringjet
    public Pair<Integer, String> GetAxisData(int row, int col, int currentAxis) {

        Pair<Integer, String> axisdata;

        if ((row == this.stone_x) && (col == this.stone_y)) {
            axisdata = new Pair<>(axes[currentAxis].getA(), axes[currentAxis].getV());
        } else {
            int currentpos = getPosition(row, col, currentAxis);
            axisdata = new Pair<>(currentpos, getSequence(row, col, currentAxis, currentpos));
        }
        return axisdata;

    }

    //SUBSTRING FINDER: megtalaljuk az elso indexet a mintanak az adott sequenceben(ha van benne)
    int Mintakereses(Pair<Integer, String> axisdata, String keyString) {

        int myPos = axisdata.getA();
        String myString = axisdata.getV();

        int keyOffset = -1;
        int keyLen = keyString.length();
        int keyPos = myString.indexOf(keyString);

        while (keyPos != -1) {
            if ((keyPos <= myPos) && (keyPos + keyLen) > myPos) {
                keyOffset = myPos - keyPos;
                return keyOffset;
            }
            keyPos = myString.indexOf(keyString, keyLen);
        }
        return keyOffset;
    }


    // Match type : Singel 1, Double 2
    //Talalunk-e egyezest a tiltott lepesekre vagy a nyertes lepesekre
    public boolean IsMatch(int row, int col, int keyType, int matchType, int restrictedAxis) {

        Pair<Integer, String> axisdata;

        int keySzam = lista.mintaListaMeret(keyType);
        String key;
        int keyOffset = -1;

        int matchSzam = 0;
        int matchLen = 0;
        int matchOffset = -1;
        for (int i = 0; i < 4; i++) {

            if (i == restrictedAxis)
                continue;
            axisdata = GetAxisData(row, col, i);

            for (int j = 0; j < keySzam; j++) {
                key = lista.getMinta(keyType, j);
                keyOffset = Mintakereses(axisdata, key);

                if (keyOffset != -1) {
                    matchSzam++;
                    if (matchSzam == matchType){
                        return true;
                    }
                    matchLen = key.length();
                    matchOffset = keyOffset;
                    restrictedAxis = i;

                    break;
                }
            }
        }
        if ((matchType == 2) && (matchSzam != 0)) {
            System.out.println("Secondmatch: Sor: " + row+ " Oszlop: "+ col+ " keyTpe: " +keyType+ " restricteda: "+ restrictedAxis+ " metchoffset: " + matchOffset+ " matchlan: " +matchLen);
            return IsSecondMatch(row, col, keyType, restrictedAxis, matchOffset, matchLen);
        }
        return false;
    }

    // Dupla harmas es negyes tesztelesenek eseteben ujra meg kell hivni ezt a fuggvenyt
    public boolean IsSecondMatch(int row, int col, int keyType, int restrictedAxis, int matchOffset, int matchLen) {

        row -= ((matchOffset) * increment[restrictedAxis][0]);
        col -= ((matchOffset) * increment[restrictedAxis][1]);
        for (int i = 0; i < matchLen; i++) {
            if ((matchOffset != i) && (tabla[row][col] == "X" ) ){
                if (IsMatch(row, col, keyType, 1, restrictedAxis)){
                    return true;
                }
            }
            row += increment[restrictedAxis][0];
            col += increment[restrictedAxis][1];
        }
        return false;
    }
    //Fajlba iras
    public void fileba(String[][] tomb, String nyertes) throws IOException {
        int lepesszam =0;
        int feherlepes =0;
        int feketelepes= 0;
        for (int i = 0; i < 17; i++) {
            for (int j = 0; j < 17; j++) {
                if(tomb[i][j].equals("X")){
                    lepesszam+= 1;
                    feketelepes += 1;
                } else if(tomb[i][j].equals("O")){
                    lepesszam += 1;
                    feherlepes += 1;
                }
            }
        }
        String out = (lepesszam +";" +feherlepes +";" +feketelepes+ ";" +nyertes);
        BufferedWriter output = new BufferedWriter(new FileWriter("eredmeny.txt",true));
        output.write(out);
        output.newLine();
        output.close();

    }

    // A tabla ertekeinek tombje. A lehelyezett kő koordinátái és annak színe.
    // + az ablakok, amiket meg kell hivni ha vege a jateknak
    public void Game(String[][] tomb, int row, int col, String ko, WinFrameWhite wfw, WinFrameBlack wfb) throws IOException {

        setTabla(tomb);
        setColor(ko);
        lista = new Minta();
        lista.Minta();
        setStone_x(row);
        setStone_y(col);
        axesinit();
        setStone(row, col);

        ///Key Types:   whiteFive 0,
        //              blackFive 1,
        //              overLine 2,
        //              inlineDouble 3,
        //              fours 4,
        //              threes 5

        if (IsMatch(stone_x, stone_y, 0, 1, -1)) {

            System.out.println("Fehernyert ,siman");
            fileba(tomb,"feher");
            wfw.run();
        }

        if (IsMatch(stone_x, stone_y, 1, 1, -1)) {

            if (IsMatch(stone_x, stone_y, 2, 1, -1)) {
                fileba(tomb,"feher");
                wfw.run();
            }

            fileba(tomb,"fekete");
            wfb.run();
        }

        if (IsMatch(stone_x, stone_y, 3, 1, -1)) {

            fileba(tomb,"feher");
            wfw.run();
        }

        if (IsMatch(stone_x, stone_y, 4, 2, -1)) {

            fileba(tomb,"feher");
            wfw.run();
        }

        if (IsMatch(stone_x, stone_y, 5, 2, -1)) {

            fileba(tomb,"feher");
            wfw.run();
        }

    }

}
