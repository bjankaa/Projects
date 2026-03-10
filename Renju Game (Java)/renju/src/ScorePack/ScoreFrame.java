package ScorePack;

import MenuPack.MenuFrame;
import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.*;
import java.util.ArrayList;
import java.util.List;

public class ScoreFrame extends JFrame {

    private ScoreData data;

    private void initComponents() throws IOException {

        this.setLayout(new BorderLayout());
        JTable tablazat = new JTable();
        JScrollPane sp = new JScrollPane(tablazat);
        this.add(sp, BorderLayout.CENTER);
        tablazat.setFillsViewportHeight(true);

        // File beolvasas
        String readLine = null;
        ScoreData tablemodel = new ScoreData();
        BufferedReader input = new BufferedReader(new FileReader("eredmeny.txt"));

        List<Score> eredmeny = new ArrayList<>();
        while((readLine = input.readLine()) != null) {

            String[] splitData = readLine.split(";");
            Score data = new Score();
            data.setLepesszam(splitData[0]);
            data.setFeherlepes(splitData[1]);
            data.setFeketelepes(splitData[2]);
            data.setNyertes(splitData[3]);
            eredmeny.add(data);
        }

        tablemodel.setList(eredmeny);
        tablazat.setModel(tablemodel);

        JPanel vissza = new JPanel();
        this.add(vissza,BorderLayout.SOUTH);
        vissza.setLayout(new FlowLayout());

        // Visszagomb
        JButton backb = new JButton("Vissza a Főoldalra");
        vissza.add(backb);
        backb.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                setVisible(false);
                MenuFrame mf = new MenuFrame();
                mf.run();
            }
        });

    }

    public ScoreFrame() throws IOException {
        super("Eredmények");
        setDefaultCloseOperation(EXIT_ON_CLOSE);
        setSize(800, 700);
        setLocationRelativeTo(null);
        setResizable(false);
        initComponents();
    }

    public void run(){
        this.setVisible(true);

    }
}
