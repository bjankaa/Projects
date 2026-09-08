package MenuPack;

import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.io.Reader;

public class RuleFrame extends JFrame {


    private void initComponents(){

        this.setLayout(new BorderLayout());
        JTextArea szoveg = new JTextArea();

        JScrollPane sp = new JScrollPane(szoveg);
        this.add(sp,BorderLayout.CENTER);

        try (BufferedReader myReader = new BufferedReader(new FileReader("Szabalyok.txt"))) {
            szoveg.read(myReader, "Szabalyok");
        } catch (IOException exp) {
            exp.printStackTrace();
        }
        szoveg.setEditable(false);

        JPanel vissza = new JPanel();
        this.add(vissza,BorderLayout.SOUTH);
        vissza.setLayout(new FlowLayout());
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



    public RuleFrame() {
        super("Szabályok");
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
