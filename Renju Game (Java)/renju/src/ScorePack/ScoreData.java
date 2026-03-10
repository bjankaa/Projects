package ScorePack;


import javax.swing.table.AbstractTableModel;
import java.util.ArrayList;
import java.util.List;

public class ScoreData extends AbstractTableModel {

    List<Score> eredmenyek = new ArrayList<>();
    String [] szoveg= {"Össz lépésszám", "Fehér lépésszám", "Fekete lépésszám", "Nyertes"};

    public void setList(List<Score> list){
        this.eredmenyek = list;
    }

    @Override
    public int getColumnCount(){
        return 4;
    }

    @Override
    public int getRowCount() {return eredmenyek.size();}

    @Override
    public Object getValueAt(int rowIndex, int columnIndex){
        Score eredmeny = eredmenyek.get(rowIndex);
        switch(columnIndex){
            case 0: return eredmeny.getLepesszam();
            case 1: return eredmeny.getFeherlepes();
            case 2: return eredmeny.getFeketelepes();
            default: return eredmeny.getNyertes();
        }
    }

    @Override
    public String getColumnName (int column){return szoveg[column];}

    @Override
    public Class getColumnClass(int columnIndex){
        if(columnIndex == 4){
            return String.class;
        } else{
            return Integer.class;
        }
    }

}
